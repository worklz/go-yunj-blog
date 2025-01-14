package shutdown

import (
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/queue"
)

// 停止队列工作任务
func StopQueueJobs() {
	queue.QueueJobScheduler.Stop()
	global.Logger.Info("队列工作任务停止调度！")
}
