package main

import (
	"go-ssip/app/common/consts"
	"go-ssip/app/service/mq/trans/initialize"
	"go-ssip/app/service/mq/trans/pkg/db"
)

func main() {
	initialize.InitLogger(consts.TransMqName)
	initialize.InitConfig()
	mysqlclient := initialize.InitDB()
	mongodbclient := initialize.InitMdb()
	consumer := initialize.InitMq()
	defer consumer.Close()

	svr := &PushMqImpl{
		DbManager: db.NewMsgManager(mysqlclient, mongodbclient),
	}
	svr.Run(consumer.Messages())
}
