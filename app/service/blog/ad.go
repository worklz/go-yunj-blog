package blog

import (
	"encoding/json"

	"github.com/worklz/yunj-blog-go/app/enum/ad/location"
	"github.com/worklz/yunj-blog-go/app/enum/redis/key"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/pkg/global"
)

type ad struct {
	Service
}

var Ad *ad

// 获取正常已发布的广告
func (a *ad) ItemsByNormalPublish(locationConst location.AdLocationConst) []model.Ad {
	var ads []model.Ad
	res, err := key.New(key.INDEX_AD_ITEMS_BY_LOCATION).SetArgs(locationConst).GetCache(0)
	if err != nil {
		global.Logger.Error("获取正常已发布的广告缓存异常！", err)
		return ads
	}

	err = json.Unmarshal(res, &ads)
	if err != nil {
		global.Logger.Error("获取正常已发布的广告缓存错误！", err)
	}
	return ads
}
