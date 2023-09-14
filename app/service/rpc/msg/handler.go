package main

import (
	"context"
	"go-ssip/app/common/kitex_gen/msg"
)

// MsgServiceImpl implements the last service interface defined in the IDL.
type MsgServiceImpl struct{}

// SendMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgReq, err error) {
	// TODO: Your code here...
	return
}

// SendGroupMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) SendGroupMsg(ctx context.Context, req *msg.SendGroupMsgReq) (resp *msg.SendGroupMsgReq, err error) {
	// TODO: Your code here...
	return
}

// GetMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) GetMsg(ctx context.Context, req *msg.GetMsgReq) (resp *msg.GetMsgResp, err error) {
	// TODO: Your code here...
	return
}

// GetGroupMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) GetGroupMsg(ctx context.Context, req *msg.GetGroupMsgReq) (resp *msg.GetGroupMsgResp, err error) {
	// TODO: Your code here...
	return
}
