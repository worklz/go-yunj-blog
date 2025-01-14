package model

type LogIp struct {
	Id         int    `gorm:"primary_key;autoIncrement;column:id" json:"id,omitempty" form:"id"`
	Ip2long    uint32 `gorm:"unique_index;column:ip2long" json:"ip2long" form:"ip2long"`
	Ip         string `gorm:"column:ip" json:"ip" form:"ip"`
	CreateTime int64  `gorm:"column:create_time" json:"create_time,omitempty" form:"create_time"`
}
