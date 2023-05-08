package organization

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/organizations", GetOrganizationList)
	g.GET("/organizations/:id", GetOrganizationByID)
	g.PUT("/organizations/:id", UpdateOrganization)
	g.POST("/organizations", NewOrganization)
	g.POST("/qrcode", GetQrCode)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/organizations/:id", WxGetOrganizationByID)
	g.POST("/wx/qrcode", WxGetQrCode)
}

func PortalRouters(g *gin.RouterGroup) {
	g.GET("/portal/organizations", PortalGetOrganizationList)
	g.GET("/portal/organizations/:id", PortalGetOrganizationByID)
}
