package initialize

import (
	"github.com/IBM/sarama"
	g "go-ssip/app/service/api/ws/global"
	"go.uber.org/zap"
	"net"
	"strconv"
)

func InitMq() sarama.PartitionConsumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	g.Logger.Info("init kafka", zap.String("addr", net.JoinHostPort(g.ServerConfig.KafkaInfo.Host, strconv.Itoa(g.ServerConfig.KafkaInfo.Port))))
	consumer, err := sarama.NewConsumer([]string{net.JoinHostPort(g.ServerConfig.KafkaInfo.Host, strconv.Itoa(g.ServerConfig.KafkaInfo.Port))}, config)
	if err != nil {
		g.Logger.Fatal("init kafka failed", zap.Error(err))
	}

	partitionConsumer, err := consumer.ConsumePartition("trans", 0, sarama.OffsetNewest)
	if err != nil {
		g.Logger.Fatal("init kafka", zap.Error(err))
	}
	return partitionConsumer
}
