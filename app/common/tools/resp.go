package tools

import (
	"errors"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/kitex_gen/base"
)

func baseResp(err errno.ErrNo) *base.BaseResponse {
	return &base.BaseResponse{
		Code: err.ErrCode,
		Msg:  err.ErrMsg,
	}
}

func BuildBaseResp(err error) *base.BaseResponse {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func ParseBaseResp(resp *base.BaseResponse) error {
	if resp.Code == errno.Success.ErrCode {
		return nil
	}
	return errno.NewErrNo(resp.Code, resp.Msg)
}
