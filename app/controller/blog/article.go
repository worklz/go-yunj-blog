package blog

import (
	"strconv"

	adLocation "github.com/worklz/yunj-blog-go/app/enum/ad/location"
	service "github.com/worklz/yunj-blog-go/app/service/blog"
	"github.com/worklz/yunj-blog-go/pkg/blog/redirect"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Controller
}

// 文章详情
func (ctrl *Article) Detail(c *gin.Context) {
	// 接收参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		redirect.Error(c, "文章异常")
		return
	}
	// 获取文章数据
	article, err := service.Article.Detail(id)
	if err != nil {
		redirect.Error(c, err.Error())
		return
	}
	// 组装渲染数据
	data := map[string]interface{}{
		"article":     article,
		"detailTopAd": service.Ad.ItemsByNormalPublish(adLocation.ARTICLE_DETAIL_TOP),
	}
	ctrl.Render(c, "article/detail", data)
}
