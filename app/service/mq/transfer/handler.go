package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"go-ssip/app/common/command"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/kitex_gen/base"
	g "go-ssip/app/service/mq/transfer/global"
	"go-ssip/app/service/mq/transfer/model"
	"go.uber.org/zap"
)

type PushMqImpl struct {
	MysqlManager
	RedisManager
	MqManager
}

type MysqlManager interface {
	Save(m []*model.Msg) error
	GetGroupMembers(g int64) ([]model.GroupMember, error)
}

type RedisManager interface {
	GetUserStatus(ctx context.Context, u int64) (string, error)
	GetAndUpdateSeq(ctx context.Context, u int64) (int64, error)
}

type MqManager interface {
	PushToPush(cmd *command.Command, st string) error
}

func (s *PushMqImpl) Run(msgs <-chan *sarama.ConsumerMessage) {
	for d := range msgs {
		var m = &base.Msg{}
		err := json.Unmarshal(d.Value, m)
		if err != nil {
			g.Logger.Error("unmarshal msg failed", zap.Error(err))
			continue
		}

		// 映射到存储的信息格式
		var ms []*model.Msg
		switch m.Type {
		case consts.MessageTypeGroupChat:
			members, err := s.MysqlManager.GetGroupMembers(m.ToGroup)
			if err != nil {
				g.Logger.Error("get group members failed", zap.Error(err))
				continue
			}
			for _, member := range members {
				m, err := s.buildMsg(context.Background(), member.UserID, m)
				if err != nil {
					g.Logger.Error("build msg failed", zap.Error(err))
					continue
				}
				ms = append(ms, m)
			}

		case consts.MessageTypeSingleChat:
			tum, err := s.buildMsg(context.Background(), m.ToUser, m)
			if err != nil {
				continue
			}
			fum, err := s.buildMsg(context.Background(), m.FromUser, m)
			if err != nil {
				continue
			}
			ms = append(ms, tum, fum)
		}

		// 持久化消息
		err = s.MysqlManager.Save(ms)
		if err != nil {
			g.Logger.Error("save msg failed", zap.Error(err))
			continue
		}
		err = s.ackClientId(context.Background(), m.FromUser, m.CliendtId)
		if err != nil {
			g.Logger.Error("send ack failed", zap.Error(err))
			// TODO 解决 ack client_id 失败问题
		}
		// 在线推送消息
		for _, m := range ms {
			if err = s.onlinePush(context.Background(), m); err != nil {
				g.Logger.Error("online push failed", zap.Error(err))
			}
		}
	}
}

func (s *PushMqImpl) buildMsg(ctx context.Context, uid int64, msg *base.Msg) (*model.Msg, error) {
	us, err := s.RedisManager.GetAndUpdateSeq(ctx, uid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			g.Logger.Info("get and update seq error", zap.Int64("user_id", uid), zap.Error(err))
			return nil, err
		}
	}
	msg.Seq = us
	content, err := json.Marshal(msg)
	if err != nil {
		g.Logger.Info("marshal msg error", zap.Error(err))
		return nil, err
	}

	var m = &model.Msg{
		UserID:  uid,
		Seq:     us,
		Content: content,
	}

	return m, nil
}

func (s *PushMqImpl) onlinePush(ctx context.Context, m *model.Msg) error {
	status, err := s.RedisManager.GetUserStatus(ctx, m.UserID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	msgBuf, _ := json.Marshal(m)
	cmd := &command.Command{
		Type:    consts.CommandTypeGetMsg,
		Payload: msgBuf,
	}

	if err = s.MqManager.PushToPush(cmd, status); err != nil {
		return err
	}
	return nil
}

func (s *PushMqImpl) ackClientId(ctx context.Context, uid int64, clientId int64) error {
	status, err := s.RedisManager.GetUserStatus(ctx, uid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	payload := &command.AckClientIdPayload{
		ClientId: clientId,
		UserId:   uid,
	}

	cmd := &command.Command{
		Type:    consts.CommandTypeAckClientId,
		Payload: payload.Encode(),
	}
	if err = s.MqManager.PushToPush(cmd, status); err != nil {
		return err
	}
	return nil
}
