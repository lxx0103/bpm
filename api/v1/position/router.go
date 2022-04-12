package position

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/positions", GetPositionList)
	g.GET("/positions/:id", GetPositionByID)
	g.PUT("/positions/:id", UpdatePosition)
	g.POST("/positions", NewPosition)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/positions", WxGetPositionList)
}
