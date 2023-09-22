package mq

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go-ssip/app/service/mq/pull/model"
)

type MsgManager struct {
	producer sarama.SyncProducer
}

func NewMsgManager(producer sarama.SyncProducer) *MsgManager {
	return &MsgManager{
		producer: producer,
	}
}

func (mm *MsgManager) PushToPush(m *model.Msg, st string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	var msg = &sarama.ProducerMessage{
		Topic: "push",
		Key:   sarama.StringEncoder(st),
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = mm.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
