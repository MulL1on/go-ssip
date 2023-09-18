package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go-ssip/app/service/mq/pull/model"
)

type MsgManager struct {
	ch *amqp.Channel
}

func NewMsgManager(ch *amqp.Channel) *MsgManager {
	return &MsgManager{
		ch: ch,
	}
}

func (mm *MsgManager) PushToPush(m *model.Msg, st string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = mm.ch.Publish(
		"",
		st,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	return nil
}
