package blog

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Controller
}

// 错误页
func (ctrl *Error) Index(c *gin.Context) {
	msg := c.Query("msg")
	msg, _ = url.PathUnescape(msg)
	ctrl.SeoAssign("错误页面")
	ctrl.Render(c, "error/index", map[string]interface{}{"msg": msg})
}
