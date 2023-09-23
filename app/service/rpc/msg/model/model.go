package model

type Msg struct {
	UserID  int64  `gorm:"column:user_id;primary_key" json:"user_id"`
	Seq     int64  `gorm:"column:seq;primary_key" json:"seq"`
	Content []byte `gorm:"column:content" json:"content"`
}

type Pr struct {
	UserID int64 `json:"user_id"`
	MinSeq int64 `json:"min_seq"`
}
