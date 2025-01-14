package service

import (
	"errors"

	"github.com/worklz/yunj-blog-go/app/enum/state"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/pkg/global"

	"github.com/jinzhu/gorm"
)

type category struct {
	Service
}

var Category *category

// 根据id获取正常的单个分类
func (c *category) GetItemNormalById(id int) model.Category {
	var category model.Category
	err := global.MySQL.Where("id = ? and state = ?", id, state.NORMAL).First(&category).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("根据id获取正常的单个分类数据异常！", err)
		}
	}
	return category
}

// 获取正常的分类数据项
func (c *category) GetItemsByNormal() []model.Category {
	var categorys []model.Category
	err := global.MySQL.
		Model(&model.Category{}).
		Where("state = ?", state.NORMAL).
		Order("sort asc, pid asc").
		Select("*").
		Find(&categorys).
		Error
	if err != nil {
		global.Logger.Error("获取正常的分类数据项异常！", err)
	}
	return categorys
}

// 获取正常的菜单树
func (c *category) GetNormalMenuTree() []model.Category {
	items := c.GetItemsByNormal()
	return c.handleNormalMenuTree(0, items)
}

// 处理正常的菜单树
func (c *category) handleNormalMenuTree(pid int, items []model.Category) []model.Category {
	menu := []model.Category{}
	for _, category := range items {
		if category.Pid != pid {
			continue
		}
		subMMenu := c.handleNormalMenuTree(category.Id, items)
		category.Sub = subMMenu
		menu = append(menu, category)
	}
	return menu
}

// 获取分类图片映射
func (c *category) GetCateImgMap(cateIds []int) map[int]string {
	res := map[int]string{}
	if len(cateIds) == 0 {
		return res
	}
	var categorys []model.Category
	err := global.MySQL.
		Model(&model.Category{}).
		Where("id in (?)", cateIds).
		Select("id,img").
		Find(&categorys).
		Error
	if err != nil {
		global.Logger.Error("获取分类数据异常！", err)
	}
	for _, category := range categorys {
		res[category.Id] = category.Img
	}
	return res
}

// 获取所有关联分类(包含自己)
func (c *category) GetRelatedCategorysById(id int, categoryMapParams ...map[int]model.Category) []model.Category {
	var categorys []model.Category
	// 创建map存储id的映射
	categoryMap := map[int]model.Category{}
	if len(categoryMapParams) <= 0 {
		allCategorys := c.GetItemsByNormal()
		for _, category := range allCategorys {
			categoryMap[category.Id] = category
		}
	} else {
		categoryMap = categoryMapParams[0]
	}
	// 获取父级分类
	if category, exists := categoryMap[id]; exists {
		categorys = append(categorys, category)
		if category.Pid > 0 {
			categorys = append(categorys, c.GetRelatedCategorysById(category.Pid, categoryMap)...)
		}
	}
	return categorys
}

// 获取默认分类
func (c *category) DefaultCategory() model.Category {
	defaultCategoryConfig := global.Config.Default.Category
	var defaultCategory model.Category
	defaultCategory.Id = defaultCategoryConfig.Id
	defaultCategory.Pid = defaultCategoryConfig.Pid
	defaultCategory.Name = defaultCategoryConfig.Name
	defaultCategory.Img = defaultCategoryConfig.Img
	return defaultCategory
}
