package api

import (
	"fmt"

	blogService "github.com/worklz/yunj-blog-go/app/service/blog"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Controller
}

func (ctrl *Article) List(c *gin.Context) {
	// 接收参数，绑定JSON数据到结构体
	var params blogService.ArticlePageListParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误！%v", err))
		return
	}
	// 验证参数
	if err := global.Validate.Struct(params); err != nil {
		response.ValidateFail(c, params, err)
		return
	}
	// 获取分页数据
	data := blogService.Article.PageList(params)
	response.Success(c, data)
}
