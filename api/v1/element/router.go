package element

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/elements", GetElementList)
	g.GET("/elements/:id", GetElementByID)
	g.PUT("/elements/:id", UpdateElement)
	g.POST("/elements", NewElement)
	g.DELETE("/elements/:id", DelElement)
}
