package config

import (
	appCorn "github.com/worklz/yunj-blog-go/app/corn"
	"github.com/worklz/yunj-blog-go/pkg/corn"
)

// 任务项集合
var TaskItems = []corn.TaskItem{
	// 每10秒执行一次
	{Enable: false, Spec: "*/10 * * * * *", Task: &appCorn.Tests{}},
	// 每天凌晨2点执行一次
	{Enable: true, Spec: "0 0 2 * * *", Task: &appCorn.BaiduSiteMap{}},
}
