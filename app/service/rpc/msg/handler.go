package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/kitex_gen/msg"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/rpc/msg/global"
	"go-ssip/app/service/rpc/msg/model"
	"go.uber.org/zap"
)

// MsgServiceImpl implements the last service interface defined in the IDL.
type MsgServiceImpl struct {
	MysqlManager
	RedisManager
	MqManager
}

type MysqlManager interface {
	GetGroupMembers(g int64)
}

type RedisManager interface {
	GetUserStatus(ctx context.Context, u int64) (string, error)
	GetAndUpdateSeq(ctx context.Context, u int64) (int64, error)
}

type MqManager interface {
	PushToTransfer(m *model.Msg) error
	PushToPush(m *model.Msg, st string) error
	PushToPull(pr *model.Pr) error
}

// SendMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	fu := req.Msg.FromUser
	switch req.Msg.Type {
	case 0:
		tu := req.Msg.ToUser
		if tu == fu {
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}

		tum, err := s.buildMsg(ctx, tu, req.Msg)
		if err != nil {
			g.Logger.Error("build tu message error")
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
		}

		fum, err := s.buildMsg(ctx, fu, req.Msg)
		if err != nil {
			g.Logger.Error("build fu message error")
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
		}

		err = s.pushMsg(ctx, tu, tum)
		if err != nil {
			g.Logger.Error("push tu message error")
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
		}

		err = s.pushMsg(ctx, fu, fum)
		if err != nil {
			g.Logger.Error("push fu message error")
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
		}

	case 1:
		//TODO group msg
	}

	return resp, nil
}

// GetMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) GetMsg(ctx context.Context, req *msg.GetMsgReq) (resp *msg.GetMsgResp, err error) {
	var pr = &model.Pr{
		UserID: req.User,
		MinSeq: req.Seq,
	}
	err = s.MqManager.PushToPull(pr)
	if err != nil {
		g.Logger.Error("push to pull error", zap.Error(err))
		resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
		return resp, nil
	}

	return resp, nil
}

func (s *MsgServiceImpl) buildMsg(ctx context.Context, u int64, msg *msg.Msg) (*model.Msg, error) {
	us, err := s.RedisManager.GetAndUpdateSeq(ctx, u)
	if err != nil {
		if err == redis.Nil {
			g.Logger.Info("get and update seq error", zap.Int64("user_id", u), zap.Error(err))
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
		UserID:  u,
		Seq:     us,
		Content: content,
	}

	return m, nil
}

func (s *MsgServiceImpl) pushMsg(ctx context.Context, u int64, m *model.Msg) error {
	err := s.MqManager.PushToTransfer(m)
	if err != nil {
		g.Logger.Error("push message to trans error", zap.Error(err))
		return err
	}
	st, err := s.RedisManager.GetUserStatus(ctx, u)
	if err != nil {
		if err == redis.Nil {
			g.Logger.Info("tu st not found")
			return nil
		}
		g.Logger.Error("get tu status error", zap.Error(err))
		return err
	}
	err = s.MqManager.PushToPush(m, st)
	if err != nil {
		g.Logger.Error("push message to push error", zap.Error(err))
		return err
	}
	return nil
}
