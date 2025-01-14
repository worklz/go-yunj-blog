package blog

import (
	"encoding/json"

	"github.com/worklz/yunj-blog-go/app/enum/redis/key"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/pkg/global"
)

type link struct {
	Service
}

var Link *link

// 获取正常已发布的连接
func (l *link) ItemsByNormalPublish() []model.Link {
	var links []model.Link
	res, err := key.New(key.INDEX_LINK_ITEMS).GetCache(0)
	if err != nil {
		global.Logger.Error("获取正常已发布的连接缓存异常！", err)
		return links
	}

	err = json.Unmarshal(res, &links)
	if err != nil {
		global.Logger.Error("获取正常已发布的连接缓存错误！", err)
	}
	return links
}
