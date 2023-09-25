package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"go-ssip/app/common/consts"
	group "go-ssip/app/common/kitex_gen/group/groupservice"
	g "go-ssip/app/service/rpc/group/global"
	"go-ssip/app/service/rpc/group/initialize"
	"go-ssip/app/service/rpc/group/pkg/gid"
	"go-ssip/app/service/rpc/group/pkg/mysql"
	"go.uber.org/zap"
	"net"
	"strconv"
)

func main() {
	initialize.InitLogger(consts.GroupServiceName)
	initialize.InitConfig()
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitRegistry(Port)
	db := initialize.InitDB()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(g.ServerConfig.Name),
		provider.WithExportEndpoint(g.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	svr := group.NewServer(&GroupServiceImpl{
		MysqlManager: mysql.NewGroupManager(db),
		IDGenerator:  gid.NewIDGenerator(),
	},

		server.WithServiceAddr(utils.NewNetAddr("tcp", net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: g.ServerConfig.Name}),
	)

	err := svr.Run()

	if err != nil {
		g.Logger.Error("run error", zap.Error(err))
	}
}
