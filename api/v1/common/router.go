package common

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/brands", GetBrandList)
	g.GET("/brands/:id", GetBrandByID)
	g.POST("/brands", NewBrand)
	g.PUT("/brands/:id", UpdateBrand)
	g.DELETE("/brands/:id", DeleteBrand)

	g.GET("/materials", GetMaterialList)
	g.GET("/materials/:id", GetMaterialByID)
	g.POST("/materials", NewMaterial)
	g.PUT("/materials/:id", UpdateMaterial)
	g.DELETE("/materials/:id", DeleteMaterial)

	g.GET("/banners", GetBannerList)
	g.GET("/banners/:id", GetBannerByID)
	g.POST("/banners", NewBanner)
	g.PUT("/banners/:id", UpdateBanner)
	g.DELETE("/banners/:id", DeleteBanner)
}

func PortalRouters(g *gin.RouterGroup) {
	g.GET("/portal/materials", PortalGetMaterialList)
	g.GET("/portal/materials/:id", PortalGetMaterialByID)
	g.GET("/portal/banners", PortalGetBannerList)
}
