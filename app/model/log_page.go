package model

type LogPage struct {
	Id     int    `gorm:"primary_key;autoIncrement;column:id" json:"id,omitempty" form:"id"`
	Md5_16 string `gorm:"unique_index;column:md5_16" json:"md5_16" form:"md5_16"`
	Url    string `gorm:"column:url" json:"url" form:"url"`
}
