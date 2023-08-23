package g

import (
	"go-ssip/app/service/rpc/user/config"
	"go.uber.org/zap"
)

var (
	GlobalConsulConfig config.ConsulConfig
	GlobalServerConfig config.ServerConfig

	Logger *zap.Logger
)
