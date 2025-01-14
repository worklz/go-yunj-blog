package model

import "github.com/worklz/yunj-blog-go/app/enum/article/status"

type Article struct {
	Model
	CategoryId        int                       `gorm:"column:category_id" json:"category_id" form:"category_id"`
	Title             string                    `gorm:"column:title" json:"title" form:"title"`
	Cover             string                    `gorm:"column:cover" json:"cover" form:"cover"`
	Keywords          string                    `gorm:"column:keywords" json:"keywords" form:"keywords"`
	ViewCount         int                       `gorm:"column:view_count" json:"view_count" form:"view_count"`
	Status            status.ArticleStatusConst `gorm:"column:status" json:"status" form:"status"`
	Attr              ArticleAttr               `gorm:"foreignKey:ArticleId;references:Id" json:"attr"`
	DisplayCreateTime string                    `gorm:"-" json:"display_create_time"` // 对外展示的创建时间
	Desc              string                    `gorm:"-" json:"desc"`                // ES需要
	Content           string                    `gorm:"-" json:"content"`
	RelatedCategorys  []Category                `gorm:"-" json:"related_categorys"`  // 关联的分类（包含当前所属分类）
	Tags              []string                  `gorm:"-" json:"tags"`               // 标签
	RecommendArticles []Article                 `gorm:"-" json:"recommend_articles"` // 推荐文章
}
