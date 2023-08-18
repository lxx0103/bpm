package upload

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/uploads", GetUploadList)
	g.POST("/uploads", NewUpload)
	g.GET("/key", GetUploadKey)
}

func WxRouters(g *gin.RouterGroup) {
	g.POST("/wx/uploads", WxNewUpload)
	g.GET("/wx/uploads", WxGetUploadList)
	g.GET("/wx/key", WxGetUploadKey)
}
