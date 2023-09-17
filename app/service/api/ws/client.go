package main

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	"github.com/spf13/cast"
	"go-ssip/app/common/kitex_gen/msg"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
	"log"
	"time"
)

type Client struct {
	id   int64
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var req = &msg.SendMsgReq{}

		err = json.Unmarshal(data, &req.Msg)
		if err != nil {
			g.Logger.Error("unmarshal message fail", zap.Error(err))
			continue
		}
		// TODO: route to different rpc service
		_, err = g.MsgClient.SendMsg(context.Background(), req)
		if err != nil {
			g.Logger.Error("msg rpc send msg error", zap.Error(err))
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetReadDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)

			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWs(_ context.Context, c *app.RequestContext) {
	id, ok := c.Get("ID")
	if !ok {
		g.Logger.Error("get id in context error")
		return
	}

	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		client := &Client{id: cast.ToInt64(id), conn: conn, send: make(chan []byte, 256)}
		hub.register <- client
		go client.writePump()
		client.readPump()
	})
	if err != nil {
		g.Logger.Error("upgrade error", zap.Error(err))
	}
}
