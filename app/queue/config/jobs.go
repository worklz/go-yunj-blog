package config

import (
	appQueue "github.com/worklz/yunj-blog-go/app/queue"
	"github.com/worklz/yunj-blog-go/pkg/queue"
)

// 队列工作项
var QueueJobs = []queue.Job{
	{Enable: true, Queue: &appQueue.Tests{}},   // 测试
	{Enable: true, Queue: &appQueue.BlogLog{}}, // 博客日志
}
