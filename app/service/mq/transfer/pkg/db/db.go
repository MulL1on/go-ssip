package db

import (
	"github.com/go-redis/redis/v8"
	"go-ssip/app/service/mq/transfer/model"
	"gorm.io/gorm"
)

type MysqlManager struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewMysqlManager(db *gorm.DB) *MysqlManager {
	return &MysqlManager{
		db: db,
	}
}

func (mm *MysqlManager) Save(ms []*model.Msg) error {
	// persist save to mysql
	tx := mm.db.Begin()
	for _, m := range ms {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
func (mm *MysqlManager) GetGroupMembers(g int64) ([]model.GroupMember, error) {
	var groupMembers []model.GroupMember
	err := mm.db.Model(&model.GroupMember{GroupID: g}).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}
	return groupMembers, nil
}
