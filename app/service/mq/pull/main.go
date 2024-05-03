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
	consumer, producer := initialize.InitMq()
	redisclient := initialize.InitRdb()
	defer consumer.Close()

	handler := &PullMqImpl{
		DbManager:    db.NewMysqlManager(mysqlclient),
		RedisManager: rdb.NewRedisManager(redisclient),
		MqManager:    mq.NewMqManager(producer),
	}

	handler.Run(consumer.Messages())
}
