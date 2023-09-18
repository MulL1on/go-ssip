package initialize

import (
	"fmt"
	"github.com/streadway/amqp"
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/mq/trans/global"
	"go.uber.org/zap"
)

func InitMq() *amqp.Connection {
	cfg := g.ServerConfig.RabbitMQInfo
	conn, err := amqp.Dial(fmt.Sprintf(consts.RabbitMqUrl, cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		g.Logger.Fatal("connect to rabbitmq failed", zap.Error(err))
	}

	return conn
}
