// Code generated by Kitex v0.7.0. DO NOT EDIT.

package groupservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	group "go-ssip/app/common/kitex_gen/group"
)

func serviceInfo() *kitex.ServiceInfo {
	return groupServiceServiceInfo
}

var groupServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "GroupService"
	handlerType := (*group.GroupService)(nil)
	methods := map[string]kitex.MethodInfo{
		"CreateGroup": kitex.NewMethodInfo(createGroupHandler, newGroupServiceCreateGroupArgs, newGroupServiceCreateGroupResult, false),
		"JoinGroup":   kitex.NewMethodInfo(joinGroupHandler, newGroupServiceJoinGroupArgs, newGroupServiceJoinGroupResult, false),
		"QuitGroup":   kitex.NewMethodInfo(quitGroupHandler, newGroupServiceQuitGroupArgs, newGroupServiceQuitGroupResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "group",
		"ServiceFilePath": "manifest/idl/rpc/group.thrift",
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

func createGroupHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*group.GroupServiceCreateGroupArgs)
	realResult := result.(*group.GroupServiceCreateGroupResult)
	success, err := handler.(group.GroupService).CreateGroup(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGroupServiceCreateGroupArgs() interface{} {
	return group.NewGroupServiceCreateGroupArgs()
}

func newGroupServiceCreateGroupResult() interface{} {
	return group.NewGroupServiceCreateGroupResult()
}

func joinGroupHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*group.GroupServiceJoinGroupArgs)
	realResult := result.(*group.GroupServiceJoinGroupResult)
	success, err := handler.(group.GroupService).JoinGroup(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGroupServiceJoinGroupArgs() interface{} {
	return group.NewGroupServiceJoinGroupArgs()
}

func newGroupServiceJoinGroupResult() interface{} {
	return group.NewGroupServiceJoinGroupResult()
}

func quitGroupHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*group.GroupServiceQuitGroupArgs)
	realResult := result.(*group.GroupServiceQuitGroupResult)
	success, err := handler.(group.GroupService).QuitGroup(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGroupServiceQuitGroupArgs() interface{} {
	return group.NewGroupServiceQuitGroupArgs()
}

func newGroupServiceQuitGroupResult() interface{} {
	return group.NewGroupServiceQuitGroupResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) CreateGroup(ctx context.Context, req *group.CreateGroupReq) (r *group.CreateGroupResp, err error) {
	var _args group.GroupServiceCreateGroupArgs
	_args.Req = req
	var _result group.GroupServiceCreateGroupResult
	if err = p.c.Call(ctx, "CreateGroup", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) JoinGroup(ctx context.Context, req *group.JoinGroupReq) (r *group.JoinGroupResp, err error) {
	var _args group.GroupServiceJoinGroupArgs
	_args.Req = req
	var _result group.GroupServiceJoinGroupResult
	if err = p.c.Call(ctx, "JoinGroup", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) QuitGroup(ctx context.Context, req *group.QuitGroupReq) (r *group.QuitGroupResp, err error) {
	var _args group.GroupServiceQuitGroupArgs
	_args.Req = req
	var _result group.GroupServiceQuitGroupResult
	if err = p.c.Call(ctx, "QuitGroup", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
