package initialize

import (
	"github.com/IBM/sarama"
	g "go-ssip/app/service/rpc/msg/global"
	"go.uber.org/zap"
	"net"
	"strconv"
	"time"
)

func InitMq() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = time.Duration(g.ServerConfig.KafkaInfo.Timeout) * time.Second
	producer, err := sarama.NewSyncProducer([]string{net.JoinHostPort(g.ServerConfig.KafkaInfo.Host, strconv.Itoa(g.ServerConfig.KafkaInfo.Port))}, config)
	if err != nil {
		g.Logger.Fatal("init kafka failed", zap.Error(err))
	}
	return producer
}
