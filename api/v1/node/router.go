package node

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/nodes", GetNodeList)
	g.GET("/nodes/:id", GetNodeByID)
	g.PUT("/nodes/:id", UpdateNode)
	g.POST("/nodes", NewNode)
	g.DELETE("/nodes/:id", DeleteNode)
}
