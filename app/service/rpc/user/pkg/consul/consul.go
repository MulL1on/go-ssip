package consul

import (
	"github.com/hashicorp/consul/api"
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/rpc/user/global"
	"net"
	"strconv"
)

type DiscoveryManager struct{}

func NewDiscoverManager() *DiscoveryManager {
	return &DiscoveryManager{}
}

func (cd *DiscoveryManager) GetWsServer() (string, error) {
	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		g.ConsulConfig.Host,
		strconv.Itoa(g.ConsulConfig.Port))

	client, err := api.NewClient(cfg)

	if err != nil {
		return "", err
	}

	services, _, err := client.Health().Service(consts.WsApiName, "", true, nil)
	if err != nil {
		return "", err
	}
	if len(services) == 0 {
		return "", nil
	}
	s := services[0]
	return net.JoinHostPort(s.Service.Address, strconv.Itoa(s.Service.Port)), nil
}
