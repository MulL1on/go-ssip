package main

import (
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/mq/trans/global"
	"go-ssip/app/service/mq/trans/initialize"
	"go-ssip/app/service/mq/trans/pkg/db"
	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger(consts.TransMqName)
	initialize.InitConfig()
	mysqlclient := initialize.InitDB()
	mongodbclient := initialize.InitMdb()
	conn := initialize.InitMq()
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		g.Logger.Fatal("declare a channel failed", zap.Error(err))
	}
	defer ch.Close()

	// 声明一个队列
	queue, err := ch.QueueDeclare(
		"trans",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		g.Logger.Fatal("declare a queue failed", zap.Error(err))
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		g.Logger.Fatal("consume msg failed", zap.Error(err))
	}

	svr := &PushMqImpl{
		DbManager: db.NewMsgManager(mysqlclient, mongodbclient),
	}
	svr.Run(msgs)
}
