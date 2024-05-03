package mq

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go-ssip/app/common/kitex_gen/base"
	"go-ssip/app/service/rpc/msg/model"
)

type MsgManager struct {
	producer sarama.SyncProducer
}

func NewMqManager(producer sarama.SyncProducer) *MsgManager {
	return &MsgManager{
		producer: producer,
	}
}

func (mm *MsgManager) PushToTransfer(m *base.Msg) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: "trans",
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = mm.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (mm *MsgManager) PushToPull(pr *model.Pr) error {
	data, err := json.Marshal(pr)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: "pull",
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = mm.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
