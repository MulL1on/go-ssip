package rpc

import (
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"go-ssip/app/common/kitex_gen/user/userservice"
	g "go-ssip/app/service/api/http/global"
	"go.uber.org/zap"
)

func initUser() {
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		g.GlobalConsulConfig.Host,
		g.GlobalConsulConfig.Port))
	if err != nil {
		g.Logger.Fatal("init consul resolver failed:", zap.Error(err))
	}
	//init OpenTelemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(g.GlobalServerConfig.UserSrvInfo.Name),
		provider.WithExportEndpoint(g.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)

	c, err := userservice.NewClient(
		g.GlobalServerConfig.UserSrvInfo.Name,
		client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()))
	client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: g.GlobalServerConfig.UserSrvInfo.Name})
	if err != nil {
		g.Logger.Fatal("init user client failed:", zap.Error(err))
	}
	g.GlobalUserClient = c
	g.Logger.Info("init user client success")
}