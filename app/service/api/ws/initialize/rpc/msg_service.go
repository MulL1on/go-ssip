package rpc

import (
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"go-ssip/app/common/kitex_gen/msg/msgservice"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
)

func initMsg() {
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		g.ConsulConfig.Host,
		g.ConsulConfig.Port))
	if err != nil {
		g.Logger.Fatal("init consul resolver failed:", zap.Error(err))
	}
	//init OpenTelemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(g.ServerConfig.MsgSrvInfo.Name),
		provider.WithExportEndpoint(g.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)

	c, err := msgservice.NewClient(
		g.ServerConfig.MsgSrvInfo.Name,
		client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()))
	client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: g.ServerConfig.MsgSrvInfo.Name})
	if err != nil {
		g.Logger.Fatal("init user client failed:", zap.Error(err))
	}
	g.MsgClient = c
	g.Logger.Info("init user client success")
}
