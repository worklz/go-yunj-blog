package es

import (
	"fmt"
	"time"

	"github.com/worklz/yunj-blog-go/pkg/global"
)

// 重置所有数据
// 会备份原索引
func Reset() error {
	if !global.Config.Elasticsearch.Enable {
		return fmt.Errorf("ES未启用！")
	}
	var err error
	var query *Query
	bkMark := time.Now().Format("20060102150405")
	global.Logger.Info(fmt.Sprintf("ES索引重置开始！%s", bkMark))
	for _, index := range indexs {
		query = NewQuery(index)
		if err = query.IndexReset(bkMark); err != nil {
			return fmt.Errorf("ES索引%s重置失败！%v", query.IndexName(), err)
		} else {
			global.Logger.Info(fmt.Sprintf("ES索引%s重置成功！%s", query.IndexName(), bkMark))
		}
	}
	// 成功日志
	global.Logger.Info(fmt.Sprintf("ES索引重置成功！%s", bkMark))
	return nil
}
