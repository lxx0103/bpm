package meeting

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/meetings", GetMeetingList)
	g.GET("/meetings/:id", GetMeetingByID)
	g.PUT("/meetings/:id", UpdateMeeting)
	g.POST("/meetings", NewMeeting)
	g.DELETE("/meetings/:id", DeleteMeeting)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/meetings", WxGetMeetingList)
}
