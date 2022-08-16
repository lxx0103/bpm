package example

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/examples", GetExampleList)
	g.GET("/examples/:id", GetExampleByID)
	g.PUT("/examples/:id", UpdateExample)
	g.POST("/examples", NewExample)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/examples", WxGetExampleList)
}
