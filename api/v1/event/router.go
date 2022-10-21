package event

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/events", GetEventList)
	g.GET("/events/:id", GetEventByID)
	g.PUT("/events/:id", UpdateEvent)
	g.GET("/checkins", GetCheckinList)
	g.GET("/events/:id/audits", GetAuditHistory)
	g.GET("/events/:id/reviews", GetReview)
	g.PUT("/events/:id/deadline", UpdateEventDeadline)
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
	g.GET("/wx/events/:id/audits", WxGetAuditHistory)
	g.POST("/wx/events/:id/reviews", WxReviewEvent)
	g.GET("/wx/events/:id/reviews", WxGetReview)
	g.PUT("/wx/events/:id/deadline", WxUpdateEventDeadline)
	g.PUT("/wx/reviews/:id/handle", WxHandleReview)
}
