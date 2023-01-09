package vendor

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/vendors", NewVendor)
	g.GET("/vendors", GetVendorList)
	g.GET("/vendors/:id", GetVendorByID)
	g.PUT("/vendors/:id", UpdateVendor)
	g.DELETE("/vendors/:id", DeleteVendor)

	// g.GET("/materials", GetMaterialList)
	// g.GET("/materials/:id", GetMaterialByID)
	// g.POST("/materials", NewMaterial)
	// g.PUT("/materials/:id", UpdateMaterial)
	// g.DELETE("/materials/:id", DeleteMaterial)
}
