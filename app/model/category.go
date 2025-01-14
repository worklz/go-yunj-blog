package model

type Category struct {
	Model
	Pid   int        `gorm:"column:pid" json:"pid" form:"pid"`
	Name  string     `gorm:"column:name" json:"name" form:"name"`
	Alias string     `gorm:"column:alias" json:"alias" form:"alias"`
	Img   string     `gorm:"column:img" json:"img" form:"img"`
	Desc  string     `gorm:"column:desc" json:"desc" form:"desc"`
	Sort  int        `gorm:"column:sort" json:"sort" form:"sort"`
	Sub   []Category // 子集
}
