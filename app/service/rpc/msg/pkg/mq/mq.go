package mq

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/spf13/cast"
	"go-ssip/app/service/rpc/msg/model"
)

type MsgManager struct {
	producer sarama.SyncProducer
}

func NewMsgManager(producer sarama.SyncProducer) *MsgManager {
	return &MsgManager{
		producer: producer,
	}
}

func (mm *MsgManager) PushToTransfer(m *model.Msg) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: "trans",
		Value: sarama.ByteEncoder(data),
		Key:   sarama.StringEncoder(cast.ToString(m.UserID)),
	}
	_, _, err = mm.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (mm *MsgManager) PushToPush(m *model.Msg, st string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: st,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = mm.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
