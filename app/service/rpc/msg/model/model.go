package model

import "time"

type Msg struct {
	UserID  int64  `json:"user_id"`
	Seq     int64  `json:"seq"`
	Content []byte `json:"content"`
}

type Pr struct {
	UserID int64 `json:"user_id"`
	MinSeq int64 `json:"min_seq"`
}

type GroupMember struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	GroupID    int64     `gorm:"column:group_id" json:"group_id"`
	UserID     int64     `gorm:"column:user_id" json:"user_id"`
	CreateTime time.Time `gorm:"create_time;autoCreateTime"`
}
