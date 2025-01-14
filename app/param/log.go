package param

import "github.com/worklz/yunj-blog-go/app/enum/log/types"

// 日志记录
type LogRecord struct {
	Guid       string             `json:"guid" validate:"required,positiveInt" message:"[guid]参数必须，且为正整数"`
	PageUrl    string             `json:"page_url" validate:"required" message:"[page_url]参数必须"`
	PageViewId string             `json:"page_view_id" validate:"required,positiveInt" message:"[page_view_id]参数必须，且为正整数"`
	Type       types.LogTypeConst `json:"type" validate:"required,oneof=11 22" message:"[type]参数需在给定范围内"`
	Referer    string             `json:"referer"`
	Title      string             `json:"title"`
	Ip         string             `json:"ip"`
	UserAgent  string             `json:"user_agent"`
	CreateTime int64              `json:"create_time"`
}
