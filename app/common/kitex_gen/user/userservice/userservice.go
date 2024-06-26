// Code generated by Kitex v0.7.0. DO NOT EDIT.

package userservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	user "go-ssip/app/common/kitex_gen/user"
)

func serviceInfo() *kitex.ServiceInfo {
	return userServiceServiceInfo
}

var userServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "UserService"
	handlerType := (*user.UserService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Auth":     kitex.NewMethodInfo(authHandler, newUserServiceAuthArgs, newUserServiceAuthResult, false),
		"Register": kitex.NewMethodInfo(registerHandler, newUserServiceRegisterArgs, newUserServiceRegisterResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "user",
		"ServiceFilePath": "manifest/idl/rpc/user.thrift",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.7.0",
		Extra:           extra,
	}
	return svcInfo
}

func authHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceAuthArgs)
	realResult := result.(*user.UserServiceAuthResult)
	success, err := handler.(user.UserService).Auth(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceAuthArgs() interface{} {
	return user.NewUserServiceAuthArgs()
}

func newUserServiceAuthResult() interface{} {
	return user.NewUserServiceAuthResult()
}

func registerHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceRegisterArgs)
	realResult := result.(*user.UserServiceRegisterResult)
	success, err := handler.(user.UserService).Register(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceRegisterArgs() interface{} {
	return user.NewUserServiceRegisterArgs()
}

func newUserServiceRegisterResult() interface{} {
	return user.NewUserServiceRegisterResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Auth(ctx context.Context, req *user.AuthReq) (r *user.AuthResp, err error) {
	var _args user.UserServiceAuthArgs
	_args.Req = req
	var _result user.UserServiceAuthResult
	if err = p.c.Call(ctx, "Auth", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Register(ctx context.Context, req *user.RegisterReq) (r *user.RegisterResp, err error) {
	var _args user.UserServiceRegisterArgs
	_args.Req = req
	var _result user.UserServiceRegisterResult
	if err = p.c.Call(ctx, "Register", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
