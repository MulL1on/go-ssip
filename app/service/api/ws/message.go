package main

import (
	"context"
	msgservice "go-ssip/app/common/kitex_gen/msg"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
)

type Message struct {
	Type     int8   `json:"type"`
	FromUser int64  `json:"from_user"`
	ToUser   int64  `json:"to_user"`
	ToGroup  int64  `json:"to_group"`
	Content  string `json:"content"`
}

func (m *Message) handle() {
	var req = &msgservice.SendMsgReq{}
	_, err := g.MsgClient.SendMsg(context.Background(), req)
	if err != nil {
		g.Logger.Error("msg service error", zap.Error(err))
	}
	//case MessageTypeImage:
	//	if c, ok := hub.clients[m.ToUser]; ok {
	//		data, err := json.Marshal(m)
	//		if err != nil {
	//			g.Logger.Error("marshal message error", zap.Error(err))
	//			return
	//		}
	//		c.conn.WriteMessage(ws.BinaryMessage, data)
	//	}
}
