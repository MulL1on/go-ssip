package initialize

import (
	"github.com/IBM/sarama"
	g "go-ssip/app/service/mq/pull/global"
	"go.uber.org/zap"
	"net"
	"strconv"
	"time"
)

func InitMq() (sarama.PartitionConsumer, sarama.SyncProducer) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = time.Duration(g.ServerConfig.KafkaInfo.Timeout) * time.Second
	consumer, err := sarama.NewConsumer([]string{net.JoinHostPort(g.ServerConfig.KafkaInfo.Host, strconv.Itoa(g.ServerConfig.KafkaInfo.Port))}, config)
	if err != nil {
		g.Logger.Fatal("init kafka consumer failed", zap.Error(err))
	}

	partitionConsumer, err := consumer.ConsumePartition("pull", 0, sarama.OffsetNewest)
	if err != nil {
		g.Logger.Fatal("init kafka consumer failed", zap.Error(err))
	}

	producer, err := sarama.NewSyncProducer([]string{net.JoinHostPort(g.ServerConfig.KafkaInfo.Host, strconv.Itoa(g.ServerConfig.KafkaInfo.Port))}, config)
	if err != nil {
		g.Logger.Fatal("init kafka failed", zap.Error(err))
	}

	return partitionConsumer, producer
}
