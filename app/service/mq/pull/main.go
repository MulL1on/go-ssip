package main

import (
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/mq/pull/global"
	"go-ssip/app/service/mq/pull/initialize"
	"go-ssip/app/service/mq/pull/pkg/db"
	"go-ssip/app/service/mq/pull/pkg/mq"
	"go-ssip/app/service/mq/pull/pkg/rdb"
	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger(consts.PullMqName)
	initialize.InitConfig()
	mysqlclient := initialize.InitDB()
	mongodbclient := initialize.InitMdb()
	conn := initialize.InitMq()
	redisclient := initialize.InitRdb()
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		g.Logger.Fatal("declare a channel failed", zap.Error(err))
	}
	defer ch.Close()

	// 声明一个队列
	queue, err := ch.QueueDeclare(
		"pull",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		g.Logger.Fatal("declare a queue failed", zap.Error(err))
	}

	prs, err := ch.Consume(
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

	handler := &PullMqImpl{
		DbManager:    db.NewMsgManager(mysqlclient, mongodbclient),
		RedisManager: rdb.NewRedisManager(redisclient),
		MqManager:    mq.NewMsgManager(ch),
	}

	handler.Run(prs)

}
