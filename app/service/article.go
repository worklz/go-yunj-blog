package service

import (
	"errors"
	"fmt"

	"github.com/worklz/yunj-blog-go/app/enum/redis/key"
	"github.com/worklz/yunj-blog-go/app/enum/state"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/util"

	"github.com/jinzhu/gorm"
)

type article struct {
	Service
}

var Article *article

// 补充封面图，没有则补充分类图片，分类图也没有，默认图
func (a *article) ItemsAppendCover(articles []model.Article) {
	if len(articles) == 0 {
		return
	}
	// 获取文章的分类ids
	cateIds := []int{}
	for _, article := range articles {
		if article.CategoryId > 0 {
			cateIds = append(cateIds, article.CategoryId)
		}
	}
	cateIds = util.IntSliceUnique(cateIds)
	// 获取分类图片映射
	cateImgMap := Category.GetCateImgMap(cateIds)

	for i, article := range articles {
		if article.Cover != "" {
			continue
		}
		if cateImg, exists := cateImgMap[article.CategoryId]; exists {
			articles[i].Cover = cateImg
		}
		if article.CategoryId == 0 {
			articles[i].Cover = global.Config.Default.Category.Img
		} else if cateImg, exists := cateImgMap[article.CategoryId]; exists {
			articles[i].Cover = cateImg
		} else {
			articles[i].Cover = global.Config.Default.Article.Cover
		}
	}
}

// 根据id获取数据
func (a *article) ItemById(id int) (model.Article, error) {
	var article model.Article
	err := global.MySQL.Where("id = ? and state <> ?", id, state.DELETED).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return article, nil
		} else {
			return article, fmt.Errorf("article数据查询失败！%v", err)
		}
	}
	return article, nil
}

// 扩展属性补充
func (a *article) ItemAttrAppend(article *model.Article) error {
	var articleAttr model.ArticleAttr
	if article.Attr.ArticleId > 0 {
		articleAttr = article.Attr
	} else {
		err := global.MySQL.Where("article_id = ?", article.Id).First(&articleAttr).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("article_attr数据查询失败！%v", err)
		}
	}
	article.Content = articleAttr.Content
	return nil
}

// 扩展属性补充
func (a *article) ItemsAttrAppend(articles []model.Article) error {
	if len(articles) == 0 {
		return nil
	}
	// 需要重新查询attr的文章ids
	ids := []int{}
	for idx, article := range articles {
		if article.Attr.ArticleId >= 0 {
			a.ItemAttrAppend(&article)
			articles[idx] = article
		} else {
			ids = append(ids, article.Id)
		}
	}
	if len(ids) <= 0 {
		return nil
	}
	var attrs []model.ArticleAttr
	err := global.MySQL.Where("article_id in (?)", ids).Find(&attrs).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("文章的article_attr数据查询失败！%v", err)
	}
	for idx, article := range articles {
		if article.Attr.ArticleId > 0 {
			continue
		}
		var attr model.ArticleAttr
		for _, articleAttr := range attrs {
			if articleAttr.ArticleId == article.Id {
				attr = articleAttr
				break
			}
		}
		article.Attr = attr
		a.ItemAttrAppend(&article)
		articles[idx] = article
	}
	return nil
}

// 浏览量++
func (a *article) ViewCountInc(article *model.Article) error {
	article.ViewCount++
	err := global.MySQL.Save(&article).Error
	if err != nil {
		return err
	}
	// 缓存数据++
	_, err = key.New(key.ARTICLE_VIEW_COUNT).IncrCache()
	if err != nil {
		global.Logger.Error("文章浏览量缓存数据自增异常！", err)
	}
	// todo 更新ES数据
	return nil
}
