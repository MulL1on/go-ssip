package initialize

import (
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/rpc/msg/global"
	"go.uber.org/zap"
	"net"
	"strconv"
)

func InitRegistry(Port int) (registry.Registry, *registry.Info) {
	r, err := consul.NewConsulRegister(net.JoinHostPort(
		g.ConsulConfig.Host,
		strconv.Itoa(g.ConsulConfig.Port)),
		consul.WithCheck(&api.AgentServiceCheck{
			Interval:                       consts.ConsulCheckInterval,
			Timeout:                        consts.ConsulCheckTimeout,
			DeregisterCriticalServiceAfter: consts.ConsulCheckDeregisterCriticalServiceAfter,
		}))

	if err != nil {
		g.Logger.Fatal("create consul register error", zap.Error(err))
	}

	// Using snowflake to generate service name.
	sf, err := snowflake.NewNode(2)
	if err != nil {
		g.Logger.Fatal("create snowflake error", zap.Error(err))
	}
	info := &registry.Info{
		ServiceName: g.ServerConfig.Name,
		Addr:        utils.NewNetAddr(consts.TCP, net.JoinHostPort(g.ServerConfig.Host, strconv.Itoa(Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
	}
	return r, info
}
