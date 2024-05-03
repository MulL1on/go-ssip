package global

import (
	"go-ssip/app/service/mq/transfer/config"
	"go.uber.org/zap"
)

var (
	ConsulConfig *config.ConsulConfig
	ServerConfig *config.ServerConfig

	Logger *zap.Logger
)
