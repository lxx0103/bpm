package upload

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/uploads", GetUploadList)
	g.POST("/uploads", NewUpload)
}

func WxRouters(g *gin.RouterGroup) {
	g.POST("/wx/uploads", NewUpload)
}
