package mq

import (
	"encoding/json"
	"github.com/IBM/sarama"
	g "go-ssip/app/service/mq/pull/global"
	"go-ssip/app/service/mq/pull/model"
	"go.uber.org/zap"
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
	g.Logger.Info("send msg", zap.Int64("user_id", m.UserID))

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	var msg = &sarama.ProducerMessage{
		Topic: st,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = mm.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
