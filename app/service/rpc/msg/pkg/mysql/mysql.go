package mysql

import (
	"gorm.io/gorm"
)

type MsgManager struct {
	db *gorm.DB
}

func NewMsgManager(db *gorm.DB) *MsgManager {
	return &MsgManager{db: db}
}

func (mm *MsgManager) GetGroupMembers(g int64) {

}
