package node

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/nodes", GetNodeList)
	g.GET("/nodes/:id", GetNodeByID)
	g.POST("/nodes", NewNode)
	g.PUT("/nodes/:id", UpdateNode)
	g.DELETE("/nodes/:id", DeleteNode)
}
