package model

type LogAttr struct {
	LogId int64  `gorm:"primary_key;column:log_id" json:"log_id,omitempty" form:"log_id"`
	Desc  string `gorm:"column:desc" json:"desc" form:"desc"`
}
