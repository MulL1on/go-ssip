package global

import (
	"go-ssip/app/service/mq/trans/config"
	"go.uber.org/zap"
)

var (
	ConsulConfig *config.ConsulConfig
	ServerConfig *config.ServerConfig

	Logger *zap.Logger
)
