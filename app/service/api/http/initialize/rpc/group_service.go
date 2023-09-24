package rpc

import (
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"go-ssip/app/common/kitex_gen/group/groupservice"
	g "go-ssip/app/service/api/http/global"
	"go.uber.org/zap"
)

func initGroup() {
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		g.ConsulConfig.Host,
		g.ConsulConfig.Port))
	if err != nil {
		g.Logger.Fatal("init consul resolver failed:", zap.Error(err))
	}
	//init OpenTelemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(g.ServerConfig.GroupSrvInfo.Name),
		provider.WithExportEndpoint(g.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)

	c, err := groupservice.NewClient(
		g.ServerConfig.GroupSrvInfo.Name,
		client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()))
	client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: g.ServerConfig.GroupSrvInfo.Name})
	if err != nil {
		g.Logger.Fatal("init group client failed:", zap.Error(err))
	}
	g.GroupClient = c
	g.Logger.Info("init group client success")
}
