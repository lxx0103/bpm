package team

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/teams", GetTeamList)
	g.GET("/teams/:id", GetTeamByID)
	g.PUT("/teams/:id", UpdateTeam)
	g.POST("/teams", NewTeam)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/teams", WxGetTeamList)
}
