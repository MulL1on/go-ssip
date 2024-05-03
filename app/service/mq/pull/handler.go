package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"go-ssip/app/common/command"
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/mq/pull/global"
	"go-ssip/app/service/mq/pull/model"
	"go.uber.org/zap"
)

type DbManager interface {
	GetMessages(u, min int64) ([]*model.Msg, error)
}

type RedisManager interface {
	GetUserStatus(ctx context.Context, u int64) (string, error)
	GetMaxSeq(ctx context.Context, u int64) (int64, error)
}

type PullMqImpl struct {
	DbManager
	RedisManager
	MqManager
}

type MqManager interface {
	PushToPush(cmd *command.Command, st string) error
}

func (s *PullMqImpl) Run(prs <-chan *sarama.ConsumerMessage) {
	for d := range prs {
		var p = &model.Pr{}
		err := json.Unmarshal(d.Value, p)
		if err != nil {
			g.Logger.Error("unmarshal msg error", zap.Error(err))
			continue
		}
		st, err := s.RedisManager.GetUserStatus(context.Background(), p.UserID)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				g.Logger.Error("u st not found", zap.Int64("user_id", p.UserID))
				continue
			}
		}

		//maxSeq, err := s.RedisManager.GetMaxSeq(context.Background(), p.UserID)

		// TODO push to push
		msgs, err := s.DbManager.GetMessages(p.UserID, p.MinSeq)
		if err != nil || len(msgs) == 0 {
			continue
		}
		for _, m := range msgs {
			data, _ := json.Marshal(m)
			cmd := &command.Command{
				Type:    consts.CommandTypeGetMsg,
				Payload: data,
			}
			err := s.MqManager.PushToPush(cmd, st)
			if err != nil {
				return
			}
		}
	}
}