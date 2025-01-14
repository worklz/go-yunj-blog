package queue

import (
	"time"

	"github.com/worklz/yunj-blog-go/app/param"
	"github.com/worklz/yunj-blog-go/app/service/blog/api/log/save"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/queue"
	"github.com/worklz/yunj-blog-go/pkg/util"
)

// 实现QueueInterface接口
type BlogLog struct {
	queue.BaseQueue
}

// 队列名称
func (bl *BlogLog) Queue() string {
	return "blog.log"
}

// 队列描述
func (bl *BlogLog) Desc() string {
	return "博客日志"
}

// 开启的协程数量
func (bl *BlogLog) CoroutineNum() int {
	return 2
}

// 每次获取消费的消息数量
func (bl *BlogLog) MessageNum() int {
	return 1
}

// 消息消费处理超时时间
func (bl *BlogLog) Timeout() time.Duration {
	return 5 * time.Second
}

// 消息消费处理函数
func (bl *BlogLog) Handler(messages []string) {
	logPrefix := queue.LogPrefix(bl)
	for _, message := range messages {
		param, err := util.JsonTo[param.LogRecord](message)
		if err != nil {
			global.Logger.WithField("queueMessage", message).WithError(err).Error(logPrefix, "消息参数错误！")
			continue
		}
		err = save.LogSave.Handler(param)
		if err != nil {
			global.Logger.WithField("queueMessage", message).WithError(err).Error(logPrefix, "消息处理异常！")
			continue
		}
	}
}
