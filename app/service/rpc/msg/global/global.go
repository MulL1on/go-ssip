package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"go-ssip/app/service/rpc/msg/config"
	"go.uber.org/zap"
)

var (
	ConsulConfig config.ConsulConfig
	ServerConfig config.ServerConfig

	Logger *zap.Logger
	Rdb    *redis.Client
	MqChan *amqp.Channel
)
