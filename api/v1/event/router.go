package event

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/events", GetEventList)
	g.GET("/events/:id", GetEventByID)
	g.PUT("/events/:id", UpdateEvent)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/events", WxGetEvents)
	g.GET("/wx/myevents", WxGetMyEvents)
}
