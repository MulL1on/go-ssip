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
	GetAndUpdateSeq(ctx context.Context, u int64) (int, error)
}

type MqManager interface {
	PushToTransfer(m *model.Msg) error
	PushToPush(m *model.Msg, st string) error
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
		fus, err := s.RedisManager.GetAndUpdateSeq(ctx, fu)
		tus, err := s.RedisManager.GetAndUpdateSeq(ctx, tu)
		m, err := json.Marshal(req.Msg)
		if err != nil {
			g.Logger.Error("marshal msg error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		var fum = &model.Msg{
			UserID:  fu,
			Seq:     fus,
			Content: m,
		}
		var tum = &model.Msg{
			UserID:  tu,
			Seq:     tus,
			Content: m,
		}
		err = s.MqManager.PushToTransfer(fum)
		if err != nil {
			g.Logger.Error("push message to trans error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		err = s.MqManager.PushToTransfer(tum)
		if err != nil {
			g.Logger.Error("push message to trans error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		st, err := s.RedisManager.GetUserStatus(ctx, tu)
		if err != nil {
			if err == redis.Nil {
				g.Logger.Info("tu st not found")
				return resp, nil
			}
			g.Logger.Error("get tu status error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		err = s.MqManager.PushToPush(tum, st)
		if err != nil {
			g.Logger.Error("push message to push error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}

		st, err = s.RedisManager.GetUserStatus(ctx, fu)
		if err != nil {
			if err == redis.Nil {
				g.Logger.Info("fu st not found")
				return resp, nil
			}
			g.Logger.Error("get fu status error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		err = s.MqManager.PushToPush(fum, st)
		if err != nil {
			g.Logger.Error("push message to push error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}

	case 1:
		//TODO group msg
	}

	return resp, nil
}

// GetMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) GetMsg(ctx context.Context, req *msg.GetMsgReq) (resp *msg.GetMsgResp, err error) {

	return
}
