package status

type ArticleStatusConst int

// 定义文章状态常量
const (
	// 草稿
	PENDING ArticleStatusConst = 11

	// 已发布
	PUBLISH ArticleStatusConst = 22
)

// 所有常量属性
var AllConstAttrs = map[ArticleStatusConst]interface{}{
	PENDING: map[string]interface{}{
		"desc": "草稿",
	},
	PUBLISH: map[string]interface{}{
		"desc": "已发布",
	},
}
