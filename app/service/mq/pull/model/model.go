package model

type Msg struct {
	ID      int64  `gorm:"column:id;primary_key" json:"id"`
	UserID  int64  `gorm:"column:user_id" json:"user_id"`
	Seq     int64  `gorm:"column:seq" json:"seq"`
	Content []byte `gorm:"content:seq" json:"content"`
}

type Pr struct {
	UserID int64 `json:"user_id"`
}
