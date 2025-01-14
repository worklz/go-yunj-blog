package record

import (
	"github.com/worklz/yunj-blog-go/app/param"
	appQueue "github.com/worklz/yunj-blog-go/app/queue"
	"github.com/worklz/yunj-blog-go/app/service/blog/api"
	"github.com/worklz/yunj-blog-go/app/service/blog/api/log"
	"github.com/worklz/yunj-blog-go/app/service/blog/api/log/save"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/queue"
)

// 日志记录
type logRecord struct {
	log.Service
}

var LogRecord *logRecord

// 执行
func (lr *logRecord) Handler(param param.LogRecord) {
	// 校验guid
	res, err := api.Guid.Check(param.Guid)
	if err != nil {
		global.Logger.WithField("param", param).WithError(err).Error("日志记录guid校验失败！")
	}
	if res != 1 {
		return
	}
	// 推送队列
	err = queue.New(&appQueue.BlogLog{}).Push(param)
	if err == nil {
		return
	}
	global.Logger.WithField("param", param).WithError(err).Error("日志记录队列推送失败！")
	// 推送队列不成功，直接保存到数据库
	err = save.LogSave.Handler(param)
	if err != nil {
		global.Logger.WithField("param", param).WithError(err).Error("日志记录保存数据库失败！", err)
	}
}
