package component

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/components", GetComponentList)
	g.GET("/components/:id", GetComponentByID)
}
