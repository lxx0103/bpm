package member

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/members", GetMemberList)
	g.POST("/members", NewMember)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/members", WxGetMemberList)
	g.POST("/wx/members", WxNewMember)
}
