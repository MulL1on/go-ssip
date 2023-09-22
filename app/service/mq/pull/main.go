package main

import (
	"go-ssip/app/common/consts"
	"go-ssip/app/service/mq/pull/initialize"
	"go-ssip/app/service/mq/pull/pkg/db"
	"go-ssip/app/service/mq/pull/pkg/mq"
	"go-ssip/app/service/mq/pull/pkg/rdb"
)

func main() {
	initialize.InitLogger(consts.PullMqName)
	initialize.InitConfig()
	mysqlclient := initialize.InitDB()
	mongodbclient := initialize.InitMdb()
	consumer, producer := initialize.InitMq()
	redisclient := initialize.InitRdb()
	defer consumer.Close()

	handler := &PullMqImpl{
		DbManager:    db.NewMsgManager(mysqlclient, mongodbclient),
		RedisManager: rdb.NewRedisManager(redisclient),
		MqManager:    mq.NewMsgManager(producer),
	}

	handler.Run(consumer.Messages())
}
