package router

import "github.com/gin-gonic/gin"

type Router func(*gin.Engine)

var (
	routers     = []Router{}
	middlewares = []gin.HandlerFunc{}
)

func Register(routes ...Router) {
	routers = append(routers, routes...)
}
func Use(handlers ...gin.HandlerFunc) {
	middlewares = append(middlewares, handlers...)
}

func InitRoutes() *gin.Engine {
	g := gin.Default()
	g.Use(middlewares...)
	// 加载路由：
	for _, route := range routers {
		route(g) // 执行自定义的路由
	}

	return g
}
