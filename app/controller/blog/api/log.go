package api

import (
	"fmt"
	"time"

	"github.com/worklz/yunj-blog-go/app/param"
	"github.com/worklz/yunj-blog-go/app/service/blog/api/log/record"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Log struct {
	Controller
}

func (ctrl *Log) Record(c *gin.Context) {
	// 接收参数，绑定JSON数据到结构体
	var param param.LogRecord
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误！%v", err))
		return
	}
	// 验证参数
	if err := global.Validate.Struct(param); err != nil {
		response.ValidateFail(c, param, err)
		return
	}
	// 补充剩余参数
	param.Ip = c.ClientIP()
	param.UserAgent = c.Request.UserAgent()
	param.CreateTime = time.Now().Unix()
	// 日志记录
	record.LogRecord.Handler(param)
	response.Success(c, nil)
}
