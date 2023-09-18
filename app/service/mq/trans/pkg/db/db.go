package db

import (
	"context"
	g "go-ssip/app/service/mq/trans/global"
	"go-ssip/app/service/mq/trans/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MsgManager struct {
	db   *gorm.DB
	coll *mongo.Collection
}

func NewMsgManager(db *gorm.DB, coll *mongo.Collection) *MsgManager {
	return &MsgManager{
		db:   db,
		coll: coll,
	}
}

func (mm *MsgManager) Save(ctx context.Context, m *model.Msg) {
	// persist save to mysql
	err := mm.db.Create(&m).Error
	if err != nil {
		g.Logger.Error("save to mysql error", zap.Error(err))
		return
	}

	// cache to mongodb
	a := bson.D{{"user_id", m.UserID}, {"seq", m.Seq}, {"content", m.Content}}
	_, err = mm.coll.InsertOne(ctx, a)
	if err != nil {
		g.Logger.Error("save to mongodb error", zap.Error(err))
		return
	}
	return
}
