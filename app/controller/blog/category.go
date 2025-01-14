package blog

import (
	"strconv"

	"github.com/worklz/yunj-blog-go/app/enum/article/sortrule"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/app/service"
	blogService "github.com/worklz/yunj-blog-go/app/service/blog"
	"github.com/worklz/yunj-blog-go/pkg/blog/redirect"

	"github.com/gin-gonic/gin"
)

type Category struct {
	Controller
}

// 分类页
func (ctrl *Category) Index(c *gin.Context) {
	// 接收参数，并获取当前分类
	idStr := c.Param("id")
	var category model.Category
	if idStr == "" {
		category = service.Category.DefaultCategory()
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			redirect.Error(c, "分类异常")
			return
		}
		category = service.Category.GetItemNormalById(id)
		if category.Id == 0 {
			redirect.Error(c, "分类数据异常")
			return
		}
	}
	// 获取文章列表数据
	params := blogService.ArticlePageListParams{
		Page:        1,
		PageSize:    8,
		Keywords:    "",
		CategoryIds: []int{category.Id},
		SortRule:    sortrule.HOT,
	}
	pageData := blogService.Article.PageList(params)
	// 组装渲染数据
	data := map[string]interface{}{
		"category":        category,
		"sortRule":        sortrule.HOT,
		"articlePageData": pageData,
	}
	ctrl.Render(c, "category/index", data)
}
