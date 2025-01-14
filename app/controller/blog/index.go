package blog

import (
	adLocation "github.com/worklz/yunj-blog-go/app/enum/ad/location"
	service "github.com/worklz/yunj-blog-go/app/service/blog"
	"github.com/worklz/yunj-blog-go/pkg/global"

	"github.com/gin-gonic/gin"
)

type Index struct {
	Controller
}

// 首页
func (ctrl *Index) Index(c *gin.Context) {
	global.Logger.Info("测试日志：", global.Config.Default.Category.Name)
	data := map[string]interface{}{
		"carousel":         service.Ad.ItemsByNormalPublish(adLocation.INDEX_CAROUSEL),
		"topCategoryItems": service.Category.GetItemsByTopNormal(),
		"articlePageData":  service.Article.PageList(),
	}
	ctrl.Render(c, "index/index", data)
}
