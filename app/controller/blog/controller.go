package blog

import (
	"fmt"
	"time"

	"github.com/worklz/yunj-blog-go/app/controller"
	adLocation "github.com/worklz/yunj-blog-go/app/enum/ad/location"
	service "github.com/worklz/yunj-blog-go/app/service/blog"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/response"
	"github.com/worklz/yunj-blog-go/pkg/util"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	controller.Controller
}

// 模板渲染
// @Param  path  string  模板路径
// @Param  data  map[string]interface{}  渲染数据
func (ctrl *Controller) Render(ctx *gin.Context, path string, data ...map[string]interface{}) {
	// 设置公共参数
	ctrl.SetCommonData()
	// 处理传入参数
	if len(data) > 0 {
		ctrl.Assign(data[0])
	}
	// 模板渲染
	response.Render(ctx, path, ctrl.Data)
}

// 设置公共参数
func (ctrl *Controller) SetCommonData() {
	// 页面浏览标记id
	pageViewId, _ := util.SnowflakeId()
	// 版权信息
	year, _ := util.ToString(time.Now().Year())
	copyrightInfo := fmt.Sprintf("<p>云静博客 © 2019-%s 云静网络 (iyunj.cn)<br><br><a href=\"https://beian.miit.gov.cn\" target=\"_blank\">渝ICP备18007862号-3</a></p>", year)
	ctrl.Assign(map[string]interface{}{
		// 站点数据
		"site": map[string]interface{}{
			"title":       global.Config.App.Name,
			"keywords":    global.Config.App.Keywords,
			"description": global.Config.App.Description,
			"version":     global.Config.App.Version,
			"author": map[string]interface{}{
				"qq":    global.Config.App.Author.QQ,
				"email": global.Config.App.Author.Email,
			},
		},
		// 页面浏览标记id
		"pageViewId": pageViewId,
		// 版权信息
		"copyrightInfo": copyrightInfo,
		// 菜单结构
		"menuLayout": service.Category.GetMenuHtmlLayout(),
		// 已发布文章数量
		"articleTotal": service.Article.TotalByNormalPublish(),
		// 文章总浏览量
		"articleViewTotal": service.Article.ViewTotal(),
		// 友情链接
		"links": service.Link.ItemsByNormalPublish(),
		// 热门文章
		"articleHotItems": service.Article.HotItems(),
		// 侧边栏广告
		"asideAd": service.Ad.ItemsByNormalPublish(adLocation.ASIDE),
	})
}

// 设置seo数据
// title: 标题
// args: 关键词，描述
func (ctrl *Controller) SeoAssign(title string, args ...string) {
	ctrl.Assign("seoTitle", title)
	if len(args) > 0 {
		ctrl.Assign("seoKeywords", args[0])
		if len(args) > 1 {
			ctrl.Assign("seoDescription", args[1])
		}
	}
}
