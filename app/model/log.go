package model

import "github.com/worklz/yunj-blog-go/app/enum/log/types"

type Log struct {
	Id         int64              `gorm:"primary_key;column:id" json:"id,omitempty" form:"id"`
	Guid       int64              `gorm:"primary_key;column:guid" json:"guid,omitempty" form:"guid"`
	PageId     int                `gorm:"column:page_id" json:"page_id" form:"page_id"`
	PageViewId int64              `gorm:"column:page_view_id" json:"page_view_id" form:"page_view_id"`
	Type       types.LogTypeConst `gorm:"column:type" json:"type" form:"type"`
	UaId       int                `gorm:"column:ua_id" json:"ua_id" form:"ua_id"`
	IpId       int                `gorm:"column:ip_id" json:"ip_id" form:"ip_id"`
	CreateTime int64              `gorm:"column:create_time" json:"create_time,omitempty" form:"create_time"`
}
