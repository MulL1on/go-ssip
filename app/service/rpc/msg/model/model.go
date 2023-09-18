package model

type Msg struct {
	UserID  int64  `gorm:"column:user_id;primary_key" json:"user_id"`
	Seq     int    `gorm:"column:seq;primary_key" json:"seq"`
	Content []byte `gorm:"column:content" json:"content"`
}
