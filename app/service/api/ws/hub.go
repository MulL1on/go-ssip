package main

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
)

var hub = newHub()

type Hub struct {
	clients     map[int64]*Client
	clientsLock sync.RWMutex

	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Register(client)
		case client := <-h.unregister:
			h.Unregister(client)
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
