package initialize

import (
	"context"
	"fmt"
	"go-ssip/app/common/consts"
	g "go-ssip/app/service/mq/pull/global"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func InitMdb() *mongo.Collection {
	c := g.ServerConfig.MongoDBInfo
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(
		fmt.Sprintf(consts.MongoURI, c.User, c.Password, c.Host, c.Port)))
	if err != nil {
		g.Logger.Fatal("connect to mongodb failed", zap.Error(err))
	}
	g.Logger.Info("init mongodb successfully")
	return client.Database(c.Name).Collection("msg")
}
