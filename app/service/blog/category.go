package blog

import (
	"fmt"

	"github.com/worklz/yunj-blog-go/app/enum/redis/key"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/app/service"
	"github.com/worklz/yunj-blog-go/pkg/global"
)

type category struct {
	Service
}

var Category *category

// 获取菜单html结构
func (c *category) GetMenuHtmlLayout() string {
	redisKey := key.New(key.INDEX_CATEGORY_MENU_HTML_LAYOUT)
	res, err := redisKey.GetCacheString(0)
	if err != nil {
		global.Logger.Error("获取菜单html结构缓存异常！", err)
		return ""
	}
	if res != "" {
		return res
	}
	// 重新获取后设置缓存
	menus := service.Category.GetNormalMenuTree()
	layout := c.handleMenuHtmlLayout(menus)
	redisKey.SetCache(layout, 0)
	return layout
}

// 处理菜单html结构
func (c *category) handleMenuHtmlLayout(menus []model.Category) string {
	layout := ""
	for _, category := range menus {
		id := category.Id
		name := category.Name
		subMenus := category.Sub
		if subMenus != nil && len(subMenus) > 0 {
			subLayout := c.handleMenuHtmlLayout(subMenus)
			layout += fmt.Sprintf(` <li class='drop-down'>
										<a href='#' title='%s'>%s</a>
										<span class='drop-down-arrow'></span>
										<ul class='drop-down-list'>
											%s
										</ul>
									</li>`, name, name, subLayout)
		} else {
			layout += fmt.Sprintf(`<li><a href='/category/%d' title='%s'>%s</a></li>`, id, name, name)
		}
	}
	return layout
}

// 获取正常的顶级分类
func (c *category) GetItemsByTopNormal() []model.Category {
	var topCategorys []model.Category
	categorys := service.Category.GetItemsByNormal()
	for _, category := range categorys {
		if category.Pid == 0 {
			topCategorys = append(topCategorys, category)
		}
	}
	return topCategorys
}
