package main

import (
	"context"
	msgservice "go-ssip/app/common/kitex_gen/msg"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
)

type message struct {
	Type     int8   `json:"type"`
	FromUser int64  `json:"from_user"`
	ToUser   int64  `json:"to_user"`
	ToGroup  int64  `json:"to_group"`
	Content  string `json:"content"`
}

func (m *message) handle() {
	var req = &msgservice.SendMsgReq{}
	_, err := g.MsgClient.SendMsg(context.Background(), req)
	if err != nil {
		g.Logger.Error("msg service error", zap.Error(err))
	}
}
