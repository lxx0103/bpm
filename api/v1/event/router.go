package event

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/events", GetEventList)
	g.GET("/events/:id", GetEventByID)
	g.PUT("/events/:id", UpdateEvent)
	// g.POST("/events", NewEvent)
	g.DELETE("/events/:id", DeleteEvent)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/events", WxGetEvents)
}
