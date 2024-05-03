package model

import "time"

type Msg struct {
	ID      int64  `gorm:"column:id;primary_key" json:"id"`
	UserID  int64  `gorm:"column:user_id" json:"user_id"`
	Seq     int64  `gorm:"column:seq" json:"seq"`
	Content []byte `gorm:"content:seq" json:"content"`
}

type GroupMember struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	GroupID    int64     `gorm:"column:group_id" json:"group_id"`
	UserID     int64     `gorm:"column:user_id" json:"user_id"`
	CreateTime time.Time `gorm:"create_time;autoCreateTime"`
}
