package api

import (
	"fmt"

	apiService "github.com/worklz/yunj-blog-go/app/service/blog/api"
	"github.com/worklz/yunj-blog-go/pkg/response"
	"github.com/worklz/yunj-blog-go/pkg/util"

	"github.com/gin-gonic/gin"
)

type Guid struct {
	Controller
}

// 校验guid有效性
// 有效返回1；临时有效返回2；无效返回0并分发一个新的10s有效的guid
func (ctrl *Guid) Check(c *gin.Context) {
	// 接收参数，绑定JSON数据到结构体
	var params apiService.GuidCheckParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误！%v", err))
		return
	}
	// 校验
	res, err := apiService.Guid.Check(params.Guid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 重新生成guid
	if res == 0 {
		res, err = apiService.Guid.Generate()
		if err != nil {
			response.Fail(c, err.Error())
			return
		}
	}
	// 转换为字符串，防止浏览器将16位以上数值四舍五入到16位
	// 如：7061389129019297560 会变为 7061389129019298000
	resStr, err := util.ToString(res)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{"guid": resStr})
}

// guid有效性反馈，客户端保存了
// 保存guid并记录ip2long
func (ctrl *Guid) Valid(c *gin.Context) {
	// 接收参数，绑定JSON数据到结构体
	var params apiService.GuidValidParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误！%v", err))
		return
	}
	// 判断是否存在
	exists, err := apiService.Guid.IsExist(params.Guid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	if exists {
		response.Success(c, nil)
		return
	}
	// 判断是否临时存在
	tempExists, err := apiService.Guid.IsTempExist(params.Guid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	if tempExists {
		err = apiService.Guid.Save(params.Guid)
		if err != nil {
			response.Fail(c, err.Error())
			return
		}
		response.Success(c, nil)
		return
	}
	response.Fail(c, "")
}
