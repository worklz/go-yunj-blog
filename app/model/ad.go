package model

import "github.com/worklz/yunj-blog-go/app/enum/ad/status"

type Ad struct {
	Model
	Location  int                  `gorm:"column:location" json:"location" form:"location"`
	Cover     string               `gorm:"column:cover" json:"name" form:"cover"`
	Link      string               `gorm:"column:link" json:"link" form:"link"`
	Desc      string               `gorm:"column:desc" json:"desc" form:"desc"`
	Sort      int                  `gorm:"column:sort" json:"sort" form:"sort"`
	ValidTime int                  `gorm:"column:valid_time" json:"valid_time" form:"valid_time"`
	Status    status.AdStatusConst `gorm:"column:status" json:"status" form:"status"`
}
