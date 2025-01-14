package router

import (
	"github.com/worklz/yunj-blog-go/app/middleware"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/html"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.App.Mode)
	router := gin.Default()
	// 中间件
	router.Use(middleware.CORS(), middleware.Log())
	// 模板引擎
	router.HTMLRender = html.NewRender("resource/view/")
	// 静态路由
	router.Static("/static", "./resource/static")
	router.Static("/upload", "./storage/upload")
	router.StaticFile("/favicon.ico", "./resource/static/favicon.ico")

	// 博客路由
	BlogRouter(router)

	return router
}
