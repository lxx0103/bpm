package organization

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/organizations", GetOrganizationList)
	g.GET("/organizations/:id", GetOrganizationByID)
	g.PUT("/organizations/:id", UpdateOrganization)
	g.POST("/organizations", NewOrganization)
}
