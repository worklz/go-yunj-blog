package index

import (
	"fmt"

	"github.com/worklz/yunj-blog-go/app/enum/state"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/app/service"
	"github.com/worklz/yunj-blog-go/pkg/global"
)

// 实现es.IndexInterface接口
type Article struct {
}

// 索引名称
func (a *Article) Name() string {
	return "article"
}

// 映射json结构
func (a *Article) Mapping() string {
	return `
		{
			"settings": {
			  "number_of_shards": 1,
			  "number_of_replicas": 0
			},
			"mappings": {
			  "properties": {
				"id": {
				  "type": "integer"
				},
				"category_id": {
				  "type": "integer"
				},
				"title": {
				  "type": "text",
				  "analyzer": "ik_max_word",
				  "search_analyzer": "ik_max_word"
				},
				"cover": {
				  "type": "keyword"
				},
				"keywords": {
				  "type": "text",
				  "analyzer": "ik_max_word",
				  "search_analyzer": "ik_max_word"
				},
				"view_count": {
				  "type": "integer"
				},
				"create_time": {
				  "type": "integer"
				},
				"last_update_time": {
				  "type": "integer"
				},
				"status": {
				  "type": "integer"
				},
				"state": {
				  "type": "integer"
				},
				"content": {
				  "type": "text",
				  "analyzer": "ik_max_word",
				  "search_analyzer": "ik_max_word"
				}
			  }
			}
		}
		`
}

// 获取要同步的数据库数据
func (a *Article) GetSyncDbDatas() (any, error) {
	var datas []model.Article
	err := global.MySQL.Preload("Attr").Where("state <> ?", state.DELETED).Find(&datas).Error
	if err != nil {
		return datas, fmt.Errorf("article同步ES数据查询失败！%v", err)
	}
	service.Article.ItemsAttrAppend(datas)
	return datas, err
}
