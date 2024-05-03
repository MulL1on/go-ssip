package main

import (
	"go-ssip/app/common/consts"
	"go-ssip/app/service/mq/transfer/initialize"
	"go-ssip/app/service/mq/transfer/pkg/db"
	"go-ssip/app/service/mq/transfer/pkg/mq"
	"go-ssip/app/service/mq/transfer/pkg/rdb"
)

func main() {
	initialize.InitLogger(consts.TransMqName)
	initialize.InitConfig()
	mysqlclient := initialize.InitDB()
	redisclient := initialize.InitRdb()
	consumer := initialize.InitMqTrans()
	producer := initialize.InitMqPush()
	defer consumer.Close()

	svr := &PushMqImpl{
		MysqlManager: db.NewMysqlManager(mysqlclient),
		RedisManager: rdb.NewRedisManager(redisclient),
		MqManager:    mq.NewMqManager(producer),
	}
	svr.Run(consumer.Messages())
}
