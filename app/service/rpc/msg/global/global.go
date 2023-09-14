package global

import (
	"go-ssip/app/service/rpc/msg/config"
	"go.uber.org/zap"
)

var (
	GlobalConsulConfig config.ConsulConfig
	GlobalServerConfig config.ServerConfig

	Logger *zap.Logger
)
