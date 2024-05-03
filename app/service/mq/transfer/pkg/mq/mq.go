package mq

import (
	"github.com/IBM/sarama"
	"go-ssip/app/common/command"
)

type MqManger struct {
	producer sarama.SyncProducer
}

func NewMqManager(producer sarama.SyncProducer) *MqManger {
	return &MqManger{
		producer: producer,
	}
}

func (mm *MqManger) PushToPush(cmd *command.Command, st string) error {
	data := cmd.Encode()
	msg := &sarama.ProducerMessage{
		Topic: st,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err := mm.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
