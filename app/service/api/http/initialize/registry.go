package initialize

import (
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/api/http/global"
	"go.uber.org/zap"
	"net"
	"strconv"
)

// InitRegistry to init consul
func InitRegistry() (registry.Registry, *registry.Info) {
	// build a consul client
	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		g.GlobalConsulConfig.Host,
		strconv.Itoa(g.GlobalConsulConfig.Port))
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		g.Logger.Fatal("init consul client failed:", zap.Error(err))
	}

	r := consul.NewConsulRegister(consulClient,
		consul.WithCheck(&api.AgentServiceCheck{
			Interval:                       consts.ConsulCheckInterval,
			Timeout:                        consts.ConsulCheckTimeout,
			DeregisterCriticalServiceAfter: consts.ConsulCheckDeregisterCriticalServiceAfter,
		}))

	// Using snowflake to generate service name.
	sf, err := snowflake.NewNode(2)
	if err != nil {
		g.Logger.Fatal("init snowflake failed:", zap.Error(err))
	}
	info := &registry.Info{
		ServiceName: g.GlobalServerConfig.Name,
		Addr: utils.NewNetAddr("tcp", net.JoinHostPort(g.GlobalServerConfig.Host,
			strconv.Itoa(g.GlobalServerConfig.Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
		Weight: registry.DefaultWeight,
	}
	return r, info
}
