package g

import (
	"go-ssip/app/service/rpc/user/config"
	"go.uber.org/zap"
)

var (
	ConsulConfig config.ConsulConfig
	ServerConfig config.ServerConfig

	Logger *zap.Logger
)
