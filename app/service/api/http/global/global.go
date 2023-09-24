package global

import (
	"go-ssip/app/common/kitex_gen/group/groupservice"
	"go-ssip/app/common/kitex_gen/user/userservice"
	"go-ssip/app/service/api/http/config"
	"go.uber.org/zap"
)

var (
	ConsulConfig *config.ConsulConfig
	ServerConfig *config.ServerConfig

	UserClient  userservice.Client
	GroupClient groupservice.Client

	Logger *zap.Logger
)
