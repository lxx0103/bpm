package project

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/projects", GetProjectList)
	g.GET("/projects/:id", GetProjectByID)
	g.PUT("/projects/:id", UpdateProject)
	g.POST("/projects", NewProject)
	g.DELETE("/projects/:id", DeleteProject)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/myprojects", WxGetMyProjects)
	g.GET("/wx/assignedprojects", WxGetAssignedProjects)

	g.GET("/wx/projects/:id", WxGetProjectByID)
	g.PUT("/wx/projects/:id", WxUpdateProject)
	g.POST("/wx/projects", WxNewProject)
	g.DELETE("/wx/projects/:id", WxDeleteProject)
}
