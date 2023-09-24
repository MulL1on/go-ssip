package mysql

import (
	"go-ssip/app/service/rpc/msg/model"
	"gorm.io/gorm"
)

type MsgManager struct {
	db *gorm.DB
}

func NewMsgManager(db *gorm.DB) *MsgManager {
	return &MsgManager{db: db}
}

func (mm *MsgManager) GetGroupMembers(g int64) ([]model.GroupMember, error) {
	var groupMembers []model.GroupMember
	err := mm.db.Model(&model.GroupMember{GroupID: g}).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}
	return groupMembers, nil
}
