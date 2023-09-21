package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/streadway/amqp"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
)

type Hub struct {
	clients     map[int64]*Client
	clientsLock sync.RWMutex
	delivery    <-chan amqp.Delivery

	register   chan *Client
	unregister chan *Client
}

type Msg struct {
	UserID  int64  `json:"user_id"`
	Seq     int64  `json:"seq"`
	Content []byte `json:"content"`
}

func newHub(delivery <-chan amqp.Delivery) *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
		delivery:   delivery,
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

func (h *Hub) Register(client *Client) {
	// TODO: pull offline messages

	// add user status
	err := g.Rdb.Set(context.Background(), cast.ToString(client.id), fmt.Sprintf(net.JoinHostPort(g.ServerConfig.Host, strconv.Itoa(g.ServerConfig.Port))), 0).Err()
	if err != nil {
		g.Logger.Error("add user redis status error", zap.Error(err))
		client.conn.Close()
		return
	}
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	h.clients[client.id] = client
}

func (h *Hub) Unregister(client *Client) {
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

func (h *Hub) Push(d amqp.Delivery) {
	var m = &Msg{}
	err := json.Unmarshal(d.Body, m)
	if err != nil {
		g.Logger.Error("unmarshal delivery error", zap.Error(err), zap.ByteString("body", d.Body))
		return
	}
	for _, c := range h.clients {
		if c.id == m.UserID {
			c.send <- m.Content
			return
		}
	}
	g.Logger.Error("no such user on this server")
}
