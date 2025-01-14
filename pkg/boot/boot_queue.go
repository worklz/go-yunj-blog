package boot

import (
	"github.com/worklz/yunj-blog-go/app/queue/config"
	"github.com/worklz/yunj-blog-go/pkg/queue"
)

// 启动队列工作
func InitQueueJobs() {
	queueJobScheduler := queue.NewQueueJobScheduler()
	queueJobScheduler.Start(config.QueueJobs)
}
