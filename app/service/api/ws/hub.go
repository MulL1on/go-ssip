package main

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/spf13/cast"
	"go-ssip/app/common/command"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/kitex_gen/msg"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Hub struct {
	clients     map[int64]*client
	clientsLock sync.RWMutex
	delivery    <-chan *sarama.ConsumerMessage

	register   chan *client
	unregister chan *client

	topic string
}

type Msg struct {
	UserID     int64  `json:"user_id"`
	Seq        int64  `json:"seq"`
	Content    []byte `json:"content"`
	ClientId   int64  `json:"client_id"`
	Timer      *time.Timer
	RetryCount int
}

func newHub(delivery <-chan *sarama.ConsumerMessage, topic string) *Hub {
	return &Hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[int64]*client),
		delivery:   delivery,

		topic: topic,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Register(client)
		case client := <-h.unregister:
			h.Unregister(client)
		case delivery := <-h.delivery:
			h.Push(delivery)
		}
	}
}

// TODO: update redis status && update redis max seq

func (h *Hub) Register(client *client) {
	// TODO: pull offline messages

	// add user status
	err := g.Rdb.Set(context.Background(), cast.ToString(client.id), h.topic, 0).Err()
	if err != nil {
		g.Logger.Error("add user redis status error", zap.Error(err))
		client.conn.Close()
		return
	}
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	h.clients[client.id] = client
	var req = &msg.GetMsgReq{
		UserId: client.id,
		Seq:    client.seq,
	}
	_, err = g.MsgClient.GetMsg(context.Background(), req)
	if err != nil {
		g.Logger.Error("sync msg error", zap.Error(err))
		client.conn.Close()
		return
	}
}

func (h *Hub) Unregister(client *client) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	// delete user status

	err := g.Rdb.Del(context.Background(), cast.ToString(client.id)).Err()
	if err != nil {
		g.Logger.Error("delete user redis status error", zap.Error(err))
		return
	}

	if _, ok := h.clients[client.id]; ok {
		delete(h.clients, client.id)
		close(client.send)
	}
}

func (h *Hub) Push(d *sarama.ConsumerMessage) {
	cmd := &command.Command{}
	cmd.Decode(d.Value)

	switch cmd.Type {
	case consts.CommandTypeGetMsg:
		var m = &Msg{}
		err := json.Unmarshal(cmd.Payload, m)
		if err != nil {
			g.Logger.Error("unmarshal delivery error", zap.Error(err), zap.ByteString("body", d.Value))
			return
		}
		h.clientsLock.RLock()
		defer h.clientsLock.RUnlock()
		c, ok := h.clients[m.UserID]
		if !ok || c == nil {
			return
		}
		if c.id == m.UserID {
			c.send <- d.Value

			// 添加一个定时器接收 ack 消息
			m.Timer = time.NewTimer(3 * time.Second)
			c.timeWheel[m.Seq] = d.Value
			go func(m *Msg) {
				for {
					select {
					case <-m.Timer.C:
						if _, ok := c.timeWheel[m.Seq]; ok {
							h.RetrySendMessage(m)
						} else {
							m.Timer.Stop()
							return
						}
					}
				}
			}(m)

			return
		}
	case consts.CommandTypeAckClientId:
		payload := command.AckClientIdPayload{}
		payload.Decode(cmd.Payload)
		c, ok := h.clients[payload.UserId]
		if !ok {
			return
		}
		c.send <- d.Value
	}
}

func (h *Hub) RetrySendMessage(m *Msg) {
	m.RetryCount++
	h.clientsLock.RLock()
	c, ok := h.clients[m.UserID]
	defer h.clientsLock.RUnlock()
	if !ok || m.RetryCount > 3 || c == nil {
		return
	}

	m.Timer.Reset(3 * time.Second)
	cmdBuf := c.timeWheel[m.Seq]
	c.send <- cmdBuf
	return
}
