package status

type AdStatusConst int

// 定义广告状态常量
const (
	// 待发布
	PENDING AdStatusConst = 11

	// 已发布
	PUBLISH AdStatusConst = 22
)

// 所有常量属性
var AllConstAttrs = map[AdStatusConst]interface{}{
	PENDING: map[string]interface{}{
		"desc": "待发布",
	},
	PUBLISH: map[string]interface{}{
		"desc": "已发布",
	},
}
