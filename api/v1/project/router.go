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
}

func WxRouters(g *gin.RouterGroup) {
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
}
