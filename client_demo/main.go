package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket" //这里使用的是 gorilla 的 websocket 库
	"go-ssip/app/common/command"
	"go-ssip/app/common/consts"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Msg struct {
	Type       int8   `json:"type"`
	FromUser   int64  `json:"from_user"`
	ToUser     int64  `json:"to_user"`
	ToGroup    int64  `json:"to_group"`
	Text       string `json:"text"`
	Seq        int64  `json:"seq"`
	ClientId   int64  `json:"client_id"`
	RetryCount int
	Timer      *time.Timer
}

const (
	zzf       int64  = 1785936564971376640
	lwj       int64  = 1705062900332236800
	lwj_token string = "v4.public.eyJhdWQiOiJ1c2VyIiwiZXhwIjoiMjAyNC0wNS0wOVQxNTo0NjowOCswODowMCIsImlhdCI6IjIwMjQtMDUtMDJUMTU6NDY6MDgrMDg6MDAiLCJpZCI6IjE3MDUwNjI5MDAzMzIyMzY4MDAiLCJpc3MiOiJnby1zc2lwIiwibmJmIjoiMjAyNC0wNS0wMlQxNTo0NjowOCswODowMCJ9ZfQq_art7iw5pp4RwTFNH0cxp9GVi-iLtOIcZm4eD97Drv6qBvOXLHdxHavQFqY-CnF-CVgAvZGTke1cJ6IKCA"
	zzf_token string = "v4.public.eyJhdWQiOiJ1c2VyIiwiZXhwIjoiMjAyNC0wNS0wOVQxNTo0OTozNCswODowMCIsImlhdCI6IjIwMjQtMDUtMDJUMTU6NDk6MzQrMDg6MDAiLCJpZCI6IjE3ODU5MzY1NjQ5NzEzNzY2NDAiLCJpc3MiOiJnby1zc2lwIiwibmJmIjoiMjAyNC0wNS0wMlQxNTo0OTozNCswODowMCJ905kns0vkxC6yUGGGPmA56VsXZFKvOiOx6mxyXo8nOuJGcdkz1dz-guQCrf6Gs6BgfpIo8scsspUoCt423aCkDA"
)

var (
	clientId  int64
	toGroup   int64
	timeWheel = make(map[int64][]byte)
)

func main() {
	args := os.Args
	token := lwj_token
	toUserId := zzf
	if len(args) > 1 {
		mode, _ := strconv.Atoi(args[1])
		switch mode {
		case 1:
			token = zzf_token
			toUserId = lwj
		case 0:
			token = lwj_token
			toUserId = zzf
		}
	}

	url := "ws://localhost:8080/"
	dialer := &websocket.Dialer{}
	header := http.Header{}
	cookie := http.Cookie{
		Name:  "seq",
		Value: "0",
	}
	header.Set("Authorization", token)
	header.Set("Cookie", cookie.String())
	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		panic(err)
	}
	conn.SetPingHandler(func(string) error {
		return conn.WriteMessage(websocket.PongMessage, nil)
	})

	c := &client{
		conn: conn,
		send: make(chan []byte, 512),
	}
	go c.read()
	go c.write()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		m := &Msg{
			Type:     consts.MessageTypeSingleChat,
			Text:     input,
			ToUser:   toUserId,
			ClientId: clientId,
		}
		mBuf, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		var cmd = &command.Command{}
		cmd.Type = consts.CommandTypeSendMsg
		cmd.Payload = mBuf

		cmdBuf := cmd.Encode()
		c.send <- cmdBuf
		m.Timer = time.NewTimer(3 * time.Second)
		timeWheel[m.ClientId] = cmdBuf
		go func(m *Msg) {
			for {
				select {
				case <-m.Timer.C:
					if _, ok := timeWheel[m.ClientId]; ok {
						if m.RetryCount > 3 {
							log.Println("exec retry count")
							m.Timer.Stop()
							return
						}
						m.RetryCount++
						c.send <- timeWheel[m.ClientId]
						m.Timer.Reset(3 * time.Second)
					}
				}
			}
		}(m)
		clientId++
	}
}

type msgPayload struct {
	ID      int64  `gorm:"column:id;primary_key" json:"id"`
	UserID  int64  `gorm:"column:user_id" json:"user_id"`
	Seq     int64  `gorm:"column:seq" json:"seq"`
	Content []byte `gorm:"content:content" json:"content"`
}

type client struct {
	conn *websocket.Conn
	send chan []byte
}

func (c *client) write() {
	for {
		select {
		case data, _ := <-c.send:
			c.conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (c *client) read() {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		cmd := &command.Command{}

		cmd.Decode(data)
		switch cmd.Type {
		case consts.CommandTypeGetMsg:
			payload := &msgPayload{}
			_ = json.Unmarshal(cmd.Payload, payload)
			fmt.Println(string(payload.Content))
			m := &Msg{}
			_ = json.Unmarshal(payload.Content, m)
			ackPayload := &command.AckMsgPayload{Seq: m.Seq}

			ackCmd := &command.Command{
				Type:    consts.CommandTypeAckMsg,
				Payload: ackPayload.Encode(),
			}
			c.send <- ackCmd.Encode()

		case consts.CommandTypeAckClientId:
			// TODO int64<-uint64
			payload := &command.AckClientIdPayload{}
			payload.Decode(cmd.Payload)
			delete(timeWheel, payload.ClientId)
		default:
			continue
		}
	}
}
