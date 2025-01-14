package sortrule

type ArticleSortRuleConst string

// 定义文章状态常量
const (
	// 时间倒序排列
	RECENT ArticleSortRuleConst = "recent"

	// 浏览量倒序排列
	HOT ArticleSortRuleConst = "hot"

	// 搜索关键词得分
	KEYWORDS_SCORE ArticleSortRuleConst = "keywords_score"
)

// 所有常量属性
var AllConstAttrs = map[ArticleSortRuleConst]interface{}{
	RECENT: map[string]interface{}{
		"desc": "时间倒序排列",
	},
	HOT: map[string]interface{}{
		"desc": "浏览量倒序排列",
	},
	KEYWORDS_SCORE: map[string]interface{}{
		"desc": "搜索关键词得分",
	},
}
