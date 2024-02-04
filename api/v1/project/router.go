package project

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/projects", GetProjectList)
	g.GET("/projects/:id", GetProjectByID)
	g.PUT("/projects/:id", UpdateProject)
	g.POST("/projects", NewProject)
	g.DELETE("/projects/:id", DeleteProject)

	g.POST("/projects/:id/reports", NewProjectReport)
	g.GET("/projects/:id/reports", GetProjectReportList)
	g.GET("/projectreports/:id", GetProjectReportByID)
	g.DELETE("/projectreports/:id", DeleteProjectReport)
	g.PUT("/projectreports/:id", UpdateProjectReport)

	g.POST("/projects/:id/records", NewProjectRecord)
	g.GET("/projects/:id/records", GetProjectRecordList)
	g.GET("/projectrecords/:id", GetProjectRecordByID)
	g.DELETE("/projectrecords/:id", DeleteProjectRecord)
	g.PUT("/projectrecords/:id", UpdateProjectRecord)

	g.GET("/projects/:id/recordStatus", GetProjectRecordStatus)
	g.GET("/projects/sumbystatus", GetProjectSumByStatus)
	g.GET("/projects/sumbyteam", GetProjectSumByTeam)
	g.GET("/projects/sumbyuser", GetProjectSumByUser)
	g.GET("/projects/sumbyarea", GetProjectSumByArea)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/projects", WxGetProjectList)
	g.GET("/wx/myprojects", WxGetMyProjects)
	g.GET("/wx/assignedprojects", WxGetAssignedProjects)

	g.GET("/wx/projects/:id", WxGetProjectByID)
	g.PUT("/wx/projects/:id", WxUpdateProject)
	g.POST("/wx/projects", WxNewProject)
	g.DELETE("/wx/projects/:id", WxDeleteProject)

	g.POST("/wx/projects/:id/reports", WxNewProjectReport)
	g.GET("/wx/projects/:id/reports", WxGetProjectReportList)
	g.GET("/wx/projectreports/:id", WxGetProjectReportByID)
	g.DELETE("/wx/projectreports/:id", WxDeleteProjectReport)
	g.PUT("/wx/projectreports/:id", WxUpdateProjectReport)
	g.POST("/wx/projectreports/:id/views", WxViewProjectReport)
	g.GET("/wx/projectreports/unread", WxGetUnreadReportList)

	g.POST("/wx/projects/:id/records", WxNewProjectRecord)
	g.GET("/wx/projects/:id/records", WxGetProjectRecordList)
	g.GET("/wx/projectrecords/:id", WxGetProjectRecordByID)
	g.DELETE("/wx/projectrecords/:id", WxDeleteProjectRecord)
	g.PUT("/wx/projectrecords/:id", WxUpdateProjectRecord)
}

func PortalRouters(g *gin.RouterGroup) {
	g.GET("/portal/projects/:id/records", PortalGetProjectRecordList)
}
