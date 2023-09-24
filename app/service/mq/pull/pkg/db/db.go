package db

import (
	g "go-ssip/app/service/mq/pull/global"
	"go-ssip/app/service/mq/pull/model"
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

func (mm *MsgManager) GetMessages(u, min int64) ([]model.Msg, error) {
	var msgs []model.Msg
	err := mm.db.Model(&model.Msg{}).Where("seq > ? and user_id = ?", min, u).Find(&msgs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			g.Logger.Error("no msg for this user", zap.Int64("user_id", u))
			return nil, gorm.ErrRecordNotFound
		}
		g.Logger.Error("find msgs in mysql error", zap.Error(err))
		return nil, err
	}
	return msgs, nil
}
