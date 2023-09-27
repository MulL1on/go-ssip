// Code generated by Kitex v0.7.0. DO NOT EDIT.

package groupservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	group "go-ssip/app/common/kitex_gen/group"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	CreateGroup(ctx context.Context, req *group.CreateGroupReq, callOptions ...callopt.Option) (r *group.CreateGroupResp, err error)
	JoinGroup(ctx context.Context, req *group.JoinGroupReq, callOptions ...callopt.Option) (r *group.JoinGroupResp, err error)
	QuitGroup(ctx context.Context, req *group.QuitGroupReq, callOptions ...callopt.Option) (r *group.QuitGroupResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kGroupServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kGroupServiceClient struct {
	*kClient
}

func (p *kGroupServiceClient) CreateGroup(ctx context.Context, req *group.CreateGroupReq, callOptions ...callopt.Option) (r *group.CreateGroupResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateGroup(ctx, req)
}

func (p *kGroupServiceClient) JoinGroup(ctx context.Context, req *group.JoinGroupReq, callOptions ...callopt.Option) (r *group.JoinGroupResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.JoinGroup(ctx, req)
}

func (p *kGroupServiceClient) QuitGroup(ctx context.Context, req *group.QuitGroupReq, callOptions ...callopt.Option) (r *group.QuitGroupResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.QuitGroup(ctx, req)
}