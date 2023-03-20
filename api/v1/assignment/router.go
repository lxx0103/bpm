package assignment

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/assignments", GetAssignmentList)
	g.GET("/assignments/:id", GetAssignmentByID)
	g.PUT("/assignments/:id", UpdateAssignment)
	g.POST("/assignments", NewAssignment)
	g.DELETE("/assignments/:id", DeleteAssignment)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/assignments", WxGetAssignmentList)
}
