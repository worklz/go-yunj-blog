package model

import "github.com/worklz/yunj-blog-go/app/enum/link/status"

type Link struct {
	Model
	Link   string                 `gorm:"column:link" json:"name" form:"link"`
	Desc   string                 `gorm:"column:desc" json:"alias" form:"desc"`
	Sort   int                    `gorm:"column:sort" json:"sort" form:"sort"`
	Status status.LinkStatusConst `gorm:"column:status" json:"status" form:"status"`
}
