package example

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/examples", GetExampleList)
	g.GET("/examples/:id", GetExampleByID)
	g.PUT("/examples/:id", UpdateExample)
	g.POST("/examples", NewExample)
	g.GET("/examples/:id/materials", GetExampleMaterialList)
	g.GET("/examples/:id/materials/:material_id", GetExampleMaterialByID)
	g.POST("/examples/:id/materials", NewExampleMaterial)
	g.PUT("/examples/:id/materials/:material_id", UpdateExampleMaterial)
	g.DELETE("/examples/:id/materials/:material_id", DeleteExampleMaterial)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/examples", WxGetExampleList)
	g.GET("/wx/examples/:id", WxGetExampleByID)
}

func PortalRouters(g *gin.RouterGroup) {
	g.GET("/portal/examples", PortalGetExampleList)
	g.GET("/portal/examples/:id", PortalGetExampleByID)
	g.GET("/portal/examples/:id/materials", PortalGetExampleMaterialList)
}
