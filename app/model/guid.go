package model

type Guid struct {
	Id         int64 `gorm:"primary_key;column:id" json:"id,omitempty" form:"id"`
	CreateTime int64 `gorm:"column:create_time" json:"create_time,omitempty" form:"create_time"`
}
