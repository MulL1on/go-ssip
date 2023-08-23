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
		"registry": kitex.NewMethodInfo(registryHandler, newUserServiceRegistryArgs, newUserServiceRegistryResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "user",
		"ServiceFilePath": "../../../../manifest/idl/rpc/user.thrift",
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

func registryHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceRegistryArgs)
	realResult := result.(*user.UserServiceRegistryResult)
	success, err := handler.(user.UserService).Registry(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceRegistryArgs() interface{} {
	return user.NewUserServiceRegistryArgs()
}

func newUserServiceRegistryResult() interface{} {
	return user.NewUserServiceRegistryResult()
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

func (p *kClient) Registry(ctx context.Context, req *user.RegistryRequest) (r *user.RegistryResponse, err error) {
	var _args user.UserServiceRegistryArgs
	_args.Req = req
	var _result user.UserServiceRegistryResult
	if err = p.c.Call(ctx, "registry", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
