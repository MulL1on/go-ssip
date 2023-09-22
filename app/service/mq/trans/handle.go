package main

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
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

func (s *PushMqImpl) Run(msgs <-chan *sarama.ConsumerMessage) {
	for d := range msgs {
		var m = &model.Msg{}
		err := json.Unmarshal(d.Value, m)
		if err != nil {
			g.Logger.Error("unmarshal msg error", zap.Error(err))
			continue
		}
		s.DbManager.Save(context.Background(), m)
	}
}
