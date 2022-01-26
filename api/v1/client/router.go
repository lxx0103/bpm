package client

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/clients", GetClientList)
	g.GET("/clients/:id", GetClientByID)
	g.PUT("/clients/:id", UpdateClient)
	g.POST("/clients", NewClient)
}
