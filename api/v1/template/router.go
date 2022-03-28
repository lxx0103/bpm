package template

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/templates", GetTemplateList)
	g.GET("/templates/:id", GetTemplateByID)
	g.PUT("/templates/:id", UpdateTemplate)
	g.POST("/templates", NewTemplate)
	g.DELETE("/templates/:id", DeleteTemplate)
}
