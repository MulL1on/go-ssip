package global

import (
	"go-ssip/app/common/kitex_gen/msg/msgservice"
	"go-ssip/app/common/kitex_gen/user/userservice"
	"go-ssip/app/service/api/http/config"
	"go.uber.org/zap"
)

var (
	ConsulConfig *config.ConsulConfig
	ServerConfig *config.ServerConfig

	UserClient userservice.Client
	MsgClient  msgservice.Client

	Logger *zap.Logger
)
