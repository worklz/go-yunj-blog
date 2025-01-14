package types

type LogTypeConst int

// 定义日志类型常量
const (
	// 页面加载
	VIEW_LOAD LogTypeConst = 11

	// 页面卸载
	VIEW_UNLOAD LogTypeConst = 22
)

// 所有常量属性
var AllConstAttrs = map[LogTypeConst]interface{}{
	VIEW_LOAD: map[string]interface{}{
		"desc": "页面加载",
	},
	VIEW_UNLOAD: map[string]interface{}{
		"desc": "页面卸载",
	},
}
