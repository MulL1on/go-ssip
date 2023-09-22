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
	"go-ssip/app/common/kitex_gen/msg/msgservice"
	g "go-ssip/app/service/rpc/msg/global"
	"go-ssip/app/service/rpc/msg/initialize"
	"go-ssip/app/service/rpc/msg/pkg/mq"
	"go-ssip/app/service/rpc/msg/pkg/mysql"
	"go-ssip/app/service/rpc/msg/pkg/rdb"
	"go.uber.org/zap"
	"net"
	"strconv"
)

func main() {
	initialize.InitLogger(consts.MsgServiceName)
	initialize.InitConfig()
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitRegistry(Port)
	redisclient := initialize.InitRdb()
	mysqlclient := initialize.InitDB()
	producer := initialize.InitMq()
	defer producer.Close()

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(g.ServerConfig.Name),
		provider.WithExportEndpoint(g.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	svr := msgservice.NewServer(&MsgServiceImpl{
		MysqlManager: mysql.NewMsgManager(mysqlclient),
		RedisManager: rdb.NewRedisManager(redisclient),
		MqManager:    mq.NewMsgManager(producer),
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
		g.Logger.Info("run error", zap.Error(err))
	}
}
