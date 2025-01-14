package model

type LogUserAgent struct {
	Id         int    `gorm:"primary_key;autoIncrement;column:id" json:"id,omitempty" form:"id"`
	Md5_16     string `gorm:"unique_index;column:md5_16" json:"md5_16" form:"md5_16"`
	Content    string `gorm:"column:content" json:"content" form:"content"`
	CreateTime int64  `gorm:"column:create_time" json:"create_time,omitempty" form:"create_time"`
}
