package client

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/clients", GetClientList)
	g.GET("/clients/:id", GetClientByID)
	g.PUT("/clients/:id", UpdateClient)
	g.POST("/clients", NewClient)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/clients", WxGetClientList)
	g.GET("/wx/clients/:id", WxGetClientByID)
	g.PUT("/wx/clients/:id", WxUpdateClient)
	g.POST("/wx/clients", WxNewClient)
}
