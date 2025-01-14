package blog

import (
	"net/url"
	"unicode/utf8"

	"github.com/worklz/yunj-blog-go/app/enum/article/sortrule"
	blogService "github.com/worklz/yunj-blog-go/app/service/blog"
	"github.com/worklz/yunj-blog-go/pkg/blog/redirect"

	"github.com/gin-gonic/gin"
)

type Search struct {
	Controller
}

// 搜索页
func (ctrl *Search) Index(c *gin.Context) {
	// 接收关键词
	keywords := c.Query("keywords")
	keywords, _ = url.PathUnescape(keywords)
	if utf8.RuneCountInString(keywords) < 2 {
		redirect.Error(c, "请输入更多的关键词信息")
		return
	}
	// 获取文章列表数据
	params := blogService.ArticlePageListParams{
		Page:        1,
		PageSize:    8,
		Keywords:    keywords,
		CategoryIds: []int{},
		SortRule:    sortrule.KEYWORDS_SCORE,
	}
	pageData := blogService.Article.PageList(params)
	// 组装渲染数据
	data := map[string]interface{}{
		"keywords":        keywords,
		"sortRule":        sortrule.KEYWORDS_SCORE,
		"articlePageData": pageData,
	}
	ctrl.Render(c, "search/index", data)
}
