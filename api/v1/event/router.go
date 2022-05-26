package event

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/events", GetEventList)
	g.GET("/events/:id", GetEventByID)
	g.PUT("/events/:id", UpdateEvent)
	g.GET("/checkins", GetCheckinList)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/events", WxGetEvents)
	g.GET("/wx/myevents", WxGetMyEvents)
	g.GET("/wx/myaudits", WxGetMyAudits)
	g.PUT("/wx/saveevents/:id", WxSaveEvent)
	g.PUT("/wx/auditevents/:id", WxAuditEvent)
	g.GET("/wx/events/:id", WxGetEventByID)
	g.PUT("/wx/events/:id", WxUpdateEvent)
	g.POST("/wx/events/:id/checkin", WxNewEventCheckin)
	g.GET("/wx/checkins", WxGetCheckinList)
}
