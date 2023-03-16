package vendors

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/vendors", NewVendors)
	g.GET("/vendors", GetVendorsList)
	g.GET("/vendors/:id", GetVendorsByID)
	g.PUT("/vendors/:id", UpdateVendors)
	g.DELETE("/vendors/:id", DeleteVendors)
}

func PortalRouters(g *gin.RouterGroup) {
	g.GET("/portal/vendors", PortalGetVendorsList)
	g.GET("/portal/vendors/:id", PortalGetVendorsByID)
}
