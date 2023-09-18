package initialize

import (
	"fmt"
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/mq/pull/global"
	"go-ssip/app/service/mq/pull/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB() *gorm.DB {
	c := g.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf(consts.MysqlDSN, c.User, c.Password, c.Host, c.Port, c.Name)

	//global mode
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		g.Logger.Fatal("connect mysql error", zap.Error(err))
	}
	db.AutoMigrate(&model.Msg{})
	g.Logger.Info("init mysql successfully")
	return db
}
