package blog

import (
	"github.com/gin-gonic/gin"
)

type AboutUs struct {
	Controller
}

// 关于我们
func (ctrl *AboutUs) Index(c *gin.Context) {
	ctrl.SeoAssign("关于我们")
	ctrl.Render(c, "about_us/index")
}
