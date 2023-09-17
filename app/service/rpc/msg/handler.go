package main

import (
	"context"
	"encoding/json"
	"go-ssip/app/common/errno"
	msg "go-ssip/app/common/kitex_gen/msg"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/rpc/msg/global"
	"go-ssip/app/service/rpc/msg/pkg/mq"
	"go-ssip/app/service/rpc/msg/pkg/rdb"
	"go.uber.org/zap"
)

// MsgServiceImpl implements the last service interface defined in the IDL.
type MsgServiceImpl struct{}

// SendMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	fu := req.Msg.FromUser
	switch req.Msg.Type {
	case 0:
		tu := req.Msg.ToUser
		fus, err := rdb.GetAndUpdateSeq(ctx, fu)
		tus, err := rdb.GetAndUpdateSeq(ctx, tu)
		m, err := json.Marshal(req.Msg)
		if err != nil {
			g.Logger.Error("marshal msg error", zap.Error(err))
		}
		var fum = &mq.Msg{
			UserID:  fu,
			Seq:     fus,
			Content: m,
		}
		var tum = &mq.Msg{
			UserID:  tu,
			Seq:     tus,
			Content: m,
		}
		err = mq.PushToMq(fum)
		if err != nil {
			g.Logger.Error("push message to mq error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		err = mq.PushToMq(tum)
		if err != nil {
			g.Logger.Error("push message to mq error", zap.Error(err))
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
	// TODO: Your code here...
	return
}
