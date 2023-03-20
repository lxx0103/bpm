package assignment

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/assignments", GetAssignmentList)
	g.GET("/assignments/:id", GetAssignmentByID)
	g.PUT("/assignments/:id", UpdateAssignment)
	g.POST("/assignments", NewAssignment)
	g.DELETE("/assignments/:id", DeleteAssignment)
	g.POST("/assignments/:id/complete", CompleteAssignment)
	g.POST("/assignments/:id/audit", AuditAssignment)
	g.GET("/assignments/my", GetMyAssignmentList)
	g.GET("/assignments/myaudit", GetMyAuditList)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/assignments", WxGetAssignmentList)
	g.GET("/wx/assignments/:id", WxGetAssignmentByID)
	g.PUT("/wx/assignments/:id", WxUpdateAssignment)
	g.POST("/wx/assignments", WxNewAssignment)
	g.DELETE("/wx/assignments/:id", WxDeleteAssignment)
	g.POST("/wx/assignments/:id/complete", WxCompleteAssignment)
	g.POST("/wx/assignments/:id/audit", WxAuditAssignment)
	g.GET("/wx/assignments/my", WxGetMyAssignmentList)
	g.GET("/wx/assignments/myaudit", WxGetMyAuditList)
}
