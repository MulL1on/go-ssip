package global

import (
	"go-ssip/app/common/kitex_gen/user/userservice"
	"go-ssip/app/service/api/http/config"
	"go.uber.org/zap"
)

var (
	GlobalConsulConfig *config.ConsulConfig
	GlobalServerConfig *config.ServerConfig

	GlobalUserClient userservice.Client

	Logger *zap.Logger
)
