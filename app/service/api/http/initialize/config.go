package initialize

import (
	"github.com/bytedance/sonic"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/api/http/global"
	"go.uber.org/zap"
	"net"
	"strconv"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile(consts.HttpApiConfigPath)
	if err := v.ReadInConfig(); err != nil {
		g.Logger.Fatal("read config failed", zap.Error(err))
	}

	if err := v.Unmarshal(&g.ConsulConfig); err != nil {
		g.Logger.Fatal("unmarshal config failed", zap.Error(err))
	}

	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		g.ConsulConfig.Host,
		strconv.Itoa(g.ConsulConfig.Port),
	)

	client, err := api.NewClient(cfg)
	if err != nil {
		g.Logger.Fatal("init consul client failed", zap.Error(err))
	}
	content, _, err := client.KV().Get(g.ConsulConfig.Key, nil)
	if err != nil {
		g.Logger.Fatal("get config failed", zap.Error(err))
	}

	err = sonic.Unmarshal(content.Value, &g.ServerConfig)
	if g.ServerConfig.Host == "" {
		g.ServerConfig.Host, err = tools.GetLocalIPv4Address()
		if err != nil {
			g.Logger.Fatal("get local ip failed", zap.Error(err))
		}
	}
}
