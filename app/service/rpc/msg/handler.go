package main

import (
	"context"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/kitex_gen/base"
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
	GetGroupMembers(g int64) ([]model.GroupMember, error)
}

type RedisManager interface {
	GetUserStatus(ctx context.Context, u int64) (string, error)
	GetAndUpdateSeq(ctx context.Context, u int64) (int64, error)
}

type MqManager interface {
	PushToTransfer(m *base.Msg) error
	PushToPull(pr *model.Pr) error
}

// SendMsg 发送消息
func (s *MsgServiceImpl) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	resp = new(msg.SendMsgResp)

	fu := req.Msg.FromUser
	switch req.Msg.Type {
	case consts.MessageTypeSingleChat:
		tu := req.Msg.ToUser
		if tu == fu || tu == 0 || fu == 0 {
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
	case consts.MessageTypeGroupChat:
		//	//TODO group msg
		tg := req.Msg.ToGroup

		if tg == 0 || fu == 0 {
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
	}
	if err = s.pushMsg(req.Msg); err != nil {
		resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
		return resp, nil
	}

	return resp, nil
}

// GetMsg 拉消息
func (s *MsgServiceImpl) GetMsg(ctx context.Context, req *msg.GetMsgReq) (resp *msg.GetMsgResp, err error) {
	resp = new(msg.GetMsgResp)
	var pr = &model.Pr{
		UserID: req.UserId,
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

// pushMsg 推送消息到 Transfer 消息队列
func (s *MsgServiceImpl) pushMsg(m *base.Msg) error {
	err := s.MqManager.PushToTransfer(m)
	if err != nil {
		g.Logger.Error("push message to trans error", zap.Error(err))
		return err
	}
	return nil
}
