package main

import (
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
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	h.clients[client.userid] = client
}

func (h *Hub) Unregister(client *Client) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	if _, ok := h.clients[client.userid]; ok {
		delete(h.clients, client.userid)
		close(client.send)
	}
}
