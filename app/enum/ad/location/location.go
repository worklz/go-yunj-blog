package location

type AdLocationConst int

// 定义广告位置常量
const (
	// 首页轮播图
	INDEX_CAROUSEL AdLocationConst = 11

	// 侧边栏
	ASIDE AdLocationConst = 22

	// 文章详情页顶部
	ARTICLE_DETAIL_TOP AdLocationConst = 33
)

// 所有常量属性
var AllConstAttrs = map[AdLocationConst]interface{}{
	INDEX_CAROUSEL: map[string]interface{}{
		"desc": "首页轮播图",
	},
	ASIDE: map[string]interface{}{
		"desc": "侧边栏",
	},
	ARTICLE_DETAIL_TOP: map[string]interface{}{
		"desc": "文章详情页顶部",
	},
}
