package initialize

import (
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/hashicorp/consul/api"
	g "go-ssip/app/service/rpc/user/global"
	"go.uber.org/zap"

	consul "github.com/kitex-contrib/registry-consul"
	"go-ssip/app/common/consts"
	"net"
	"strconv"
)

func InitRegistry(Port int) (registry.Registry, *registry.Info) {
	r, err := consul.NewConsulRegister(net.JoinHostPort(
		g.GlobalConsulConfig.Host,
		strconv.Itoa(g.GlobalConsulConfig.Port)),
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
		ServiceName: g.GlobalServerConfig.Name,
		Addr:        utils.NewNetAddr(consts.TCP, net.JoinHostPort(g.GlobalServerConfig.Host, strconv.Itoa(Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
	}
	return r, info
}