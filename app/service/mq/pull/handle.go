package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	g "go-ssip/app/service/mq/pull/global"
	"go-ssip/app/service/mq/pull/model"
	"go.uber.org/zap"
)

type DbManager interface {
	GetMessages(u int64) ([]model.Msg, error)
}

type RedisManager interface {
	GetUserStatus(ctx context.Context, u int64) (string, error)
}

type PullMqImpl struct {
	DbManager
	RedisManager
	MqManager
}

type MqManager interface {
	PushToPush(m *model.Msg, st string) error
}

func (s *PullMqImpl) Run(prs <-chan amqp.Delivery) {
	for d := range prs {
		var p = &model.Pr{}
		err := json.Unmarshal(d.Body, p)
		if err != nil {
			g.Logger.Error("unmarshal msg error", zap.Error(err))
			continue
		}
		st, err := s.RedisManager.GetUserStatus(context.Background(), p.UserID)
		if err != nil {
			if err == redis.Nil {
				g.Logger.Error("u st not found", zap.Int64("user_id", p.UserID))
				continue
			}
		}

		// TODO push to push

	}
}
