package initialize

import (
	"github.com/bytedance/sonic"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/rpc/group/global"
	"go.uber.org/zap"
	"net"
	"strconv"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile(consts.GroupSrvConfigPath)
	if err := v.ReadInConfig(); err != nil {
		g.Logger.Fatal("read config file failed", zap.Error(err))
	}
	if err := v.Unmarshal(&g.ConsulConfig); err != nil {
		g.Logger.Fatal("unmarshal config file failed", zap.Error(err))
	}
	g.Logger.Info("read config file successfully", zap.Any("config", g.ConsulConfig))

	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		g.ConsulConfig.Host,
		strconv.Itoa(g.ConsulConfig.Port))
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		g.Logger.Fatal("create consul client failed", zap.Error(err))
	}
	content, _, err := consulClient.KV().Get(g.ConsulConfig.Key, nil)
	if err != nil {
		g.Logger.Fatal("get config from consul failed", zap.Error(err))
	}
	err = sonic.Unmarshal(content.Value, &g.ServerConfig)
	if err != nil {
		g.Logger.Fatal("unmarshal config from consul failed", zap.Error(err))
	}
	g.Logger.Info("get config from consul successfully", zap.Any("config", g.ServerConfig))

	if g.ServerConfig.Host == "" {
		g.ServerConfig.Host, err = tools.GetLocalIPv4Address()
		if err != nil {
			g.Logger.Fatal("get local ip address failed", zap.Error(err))
		}
	}
}
