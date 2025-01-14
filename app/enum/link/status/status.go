package status

type LinkStatusConst int

// 定义友情链接状态常量
const (
	// 待发布
	PENDING LinkStatusConst = 11

	// 已发布
	PUBLISH LinkStatusConst = 22
)

// 所有常量属性
var AllConstAttrs = map[LinkStatusConst]interface{}{
	PENDING: map[string]interface{}{
		"desc": "待发布",
	},
	PUBLISH: map[string]interface{}{
		"desc": "已发布",
	},
}
