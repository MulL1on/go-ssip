package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/kitex_gen/msg"
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
		if tu == fu {
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		fus, err := rdb.GetAndUpdateSeq(ctx, fu)
		tus, err := rdb.GetAndUpdateSeq(ctx, tu)
		m, err := json.Marshal(req.Msg)
		if err != nil {
			g.Logger.Error("marshal msg error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
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
		err = mq.PushToTransfer(fum)
		if err != nil {
			g.Logger.Error("push message to trans error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		err = mq.PushToTransfer(tum)
		if err != nil {
			g.Logger.Error("push message to trans error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		st, err := rdb.GetUserStatus(ctx, tu)
		if err != nil {
			if err == redis.Nil {
				g.Logger.Info("tu st not found")
				return resp, nil
			}
			g.Logger.Error("get tu status error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		err = mq.PushToPush(tum, st)
		if err != nil {
			g.Logger.Error("push message to push error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}

		st, err = rdb.GetUserStatus(ctx, fu)
		if err != nil {
			if err == redis.Nil {
				g.Logger.Info("fu st not found")
				return resp, nil
			}
			g.Logger.Error("get fu status error", zap.Error(err))
			resp.BaseResp = tools.BuildBaseResp(errno.MsgSrvErr)
			return resp, nil
		}
		err = mq.PushToPush(fum, st)
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
	// TODO: Your code here...
	return
}
