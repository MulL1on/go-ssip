package main

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	g "go-ssip/app/service/mq/trans/global"
	"go-ssip/app/service/mq/trans/model"
	"go.uber.org/zap"
)

type PushMqImpl struct {
	DbManager
}

type DbManager interface {
	Save(ctx context.Context, m *model.Msg)
}

func (s *PushMqImpl) Run(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var m = &model.Msg{}
		err := json.Unmarshal(d.Body, m)
		if err != nil {
			g.Logger.Error("unmarshal msg error", zap.Error(err))
			continue
		}
		s.DbManager.Save(context.Background(), m)
	}
}
