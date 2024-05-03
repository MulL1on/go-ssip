package main

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	"github.com/spf13/cast"
	"go-ssip/app/common/command"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/kitex_gen/base"
	"go-ssip/app/common/kitex_gen/msg"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

type client struct {
	clientId  int64
	id        int64
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	seq       int64
	timeWheel map[int64][]byte
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

func (c *client) readPump() {
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

		cmd := &command.Command{}
		cmd.Decode(data)

		switch cmd.Type {

		case consts.CommandTypeAckMsg:
			// 获取 seq
			payload := &command.AckMsgPayload{}
			payload.Decode(cmd.Payload)
			m, ok := c.timeWheel[payload.Seq]
			if !ok || m == nil {
				continue
			}
			// 加锁保护
			delete(c.timeWheel, payload.Seq)
		case consts.CommandTypeSendMsg:
			var m = &base.Msg{}

			err = json.Unmarshal(cmd.Payload, &m)
			if err != nil {
				g.Logger.Error("unmarshal message fail", zap.Error(err))
				continue
			}

			// 判断 client id 是否相同，否则抛弃
			if m.ClientId != c.clientId {
				continue
			}
			c.clientId++

			// 拿到客户端的用户 id
			m.FromUser = c.id

			// 构建 rpc 请求
			var req = &msg.SendMsgReq{}
			req.Msg = m
			_, err = g.MsgClient.SendMsg(context.Background(), req)
			if err != nil {
				g.Logger.Error("msg rpc send msg error", zap.Error(err))
			}
		default:
			g.Logger.Error("unhandled command type")
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case m, ok := <-c.send:
			c.conn.SetReadDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(m)

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
	seqByte := c.Cookie("seq")
	if len(seqByte) == 0 {
		c.JSON(http.StatusBadRequest, "seq not set")
		return
	}
	seq, err := strconv.ParseInt(string(seqByte), 10, 64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal error")
		return
	}

	id, ok := c.Get("ID")
	if !ok {
		g.Logger.Error("get id in context error")
		return
	}

	err = upgrader.Upgrade(c, func(conn *websocket.Conn) {
		cli := &client{id: cast.ToInt64(id), conn: conn, send: make(chan []byte, 256), seq: cast.ToInt64(seq), timeWheel: make(map[int64][]byte)}
		hub.register <- cli
		go cli.writePump()
		cli.readPump()
	})
	if err != nil {
		g.Logger.Error("upgrade error", zap.Error(err))
	}
}
