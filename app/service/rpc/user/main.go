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
	"go-ssip/app/common/kitex_gen/user/userservice"
	g "go-ssip/app/service/rpc/user/global"
	"go-ssip/app/service/rpc/user/initialize"
	"go-ssip/app/service/rpc/user/pkg/md5"
	"go-ssip/app/service/rpc/user/pkg/mysql"
	"go-ssip/app/service/rpc/user/pkg/paseto"
	"go-ssip/app/service/rpc/user/pkg/uid"
	"go.uber.org/zap"
	"log"
	"net"
	"strconv"
)

func main() {
	initialize.InitLogger(consts.UserServiceName)
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

	//TODO: 生成自己的pubkey和prikey
	tg, err := paseto.NewTokenGenerator()

	if err != nil {
		g.Logger.Fatal("create token generator error", zap.Error(err))
	}

	svr := userservice.NewServer(&UserServiceImpl{
		EncryptManager: &md5.EncryptManager{Salt: g.ServerConfig.MysqlInfo.Salt},
		MysqlManager:   mysql.NewUserManager(db, g.ServerConfig.MysqlInfo.Salt),
		IDGenerator:    uid.NewIDGenerator(),
		TokenGenerator: tg,
	},

		server.WithServiceAddr(utils.NewNetAddr("tcp", net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: g.ServerConfig.Name}),
	)

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
