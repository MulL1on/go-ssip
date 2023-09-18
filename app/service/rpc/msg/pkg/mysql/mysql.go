package mysql

import (
	"gorm.io/gorm"
)

type MsgManager struct {
	db *gorm.DB
}

type Msg struct {
	ID      int64  `gorm:"column:id;primary_key" json:"id"`
	UserID  int64  `gorm:"column:user_id" json:"user_id"`
	Seq     int    `gorm:"column:seq" json:"seq"`
	Content []byte `gorm:"content:seq" json:"content"`
}

func NewMsgManager(db *gorm.DB) *MsgManager {
	return &MsgManager{db: db}
}

func (mm *MsgManager) GetGroupMembers(g int64) {

}
