package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	g "go-ssip/app/service/rpc/msg/global"
)

type Msg struct {
	UserID  int64  `gorm:"column:user_id;primary_key" json:"user_id"`
	Seq     int    `gorm:"column:seq;primary_key" json:"seq"`
	Content []byte `gorm:"column:content" json:"content"`
}

func PushToTransfer(m *Msg) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = g.MqChan.Publish(
		"",
		"trans",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	return nil
}

func PushToPush(m *Msg, st string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = g.MqChan.Publish(
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
