package global

import (
	"go-ssip/app/common/kitex_gen/msg/msgservice"
	"go-ssip/app/service/api/ws/config"
	"go.uber.org/zap"
)

var (
	ConsulConfig *config.ConsulConfig
	ServerConfig *config.ServerConfig

	MsgClient msgservice.Client

	Logger *zap.Logger
)
