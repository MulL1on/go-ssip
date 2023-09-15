package main

import (
	"context"
	msg "go-ssip/app/common/kitex_gen/msg"
)

// MsgServiceImpl implements the last service interface defined in the IDL.
type MsgServiceImpl struct{}

// SendMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	// TODO: Your code here...
	return
}

// GetMsg implements the MsgServiceImpl interface.
func (s *MsgServiceImpl) GetMsg(ctx context.Context, req *msg.GetMsgReq) (resp *msg.GetMsgResp, err error) {
	// TODO: Your code here...
	return
}
