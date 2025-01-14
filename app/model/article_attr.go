package model

type ArticleAttr struct {
	ArticleId int    `gorm:"primary_key;column:article_id" json:"article_id,omitempty" form:"article_id"`
	Content   string `gorm:"column:content" json:"content" form:"content"`
}
