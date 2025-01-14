package redirect

import (
	"net/url"

	"github.com/worklz/yunj-blog-go/pkg/redirect"

	"github.com/gin-gonic/gin"
)

// 重定向到错误页面
func Error(c *gin.Context, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	url := "/error?msg=" + url.PathEscape(message)
	redirect.To(c, url)
}
