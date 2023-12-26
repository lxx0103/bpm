package project

import (
	"bpm/core/response"
	"bpm/service"
	"errors"

	"github.com/gin-gonic/gin"
)

// @Summary 项目列表
// @Id M001
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "项目名称"
// @Param type query int false "项目类型"
// @Success 200 object response.ListRes{data=[]ProjectResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects [GET]
func GetProjectList(c *gin.Context) {
	var filter ProjectFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := projectService.GetProjectList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建项目
// @Id M002
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_info body ProjectNew true "项目信息"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects [POST]
func NewProject(c *gin.Context) {
	var project ProjectNew
	if err := c.ShouldBindJSON(&project); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	project.User = claims.Username
	project.UserID = claims.UserID
	organizationID := claims.OrganizationID
	projectService := NewProjectService()
	new, err := projectService.NewProject(project, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取项目
// @Id M003
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id [GET]
func GetProjectByID(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	projectService := NewProjectService()
	project, err := projectService.GetProjectByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, project)

}

// @Summary 根据ID更新项目
// @Id M004
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Param project_info body ProjectUpdate true "项目信息"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id [PUT]
func UpdateProject(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var project ProjectUpdate
	if err := c.ShouldBindJSON(&project); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	project.User = claims.Username
	organizationID := claims.OrganizationID
	projectService := NewProjectService()
	new, err := projectService.UpdateProject(uri.ID, project, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除项目
// @Id M005
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id [DELETE]
func DeleteProject(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	projectService := NewProjectService()
	err := projectService.DeleteProject(uri.ID, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 获取我创建的项目
// @Id M006
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param status query string true "显示所有all/激活active"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param type query int false "项目类型"
// @Success 200 object response.ListRes{data=[]Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/myprojects [GET]
func WxGetMyProjects(c *gin.Context) {
	var filter MyProjectFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	var count int
	var list *[]ProjectResponse
	if claims.UserType == 2 {
		count, list, err = projectService.GetMyProject(filter, claims.Username, claims.OrganizationID)
		if err != nil {
			response.ResponseError(c, "DatabaseError", err)
			return
		}
	} else if claims.UserType == 3 {
		count, list, err = projectService.GetClientProject(filter, claims.UserID, claims.OrganizationID)
		if err != nil {
			response.ResponseError(c, "DatabaseError", err)
			return
		}
	} else {
		response.ResponseError(c, "用户类型错误", errors.New("管理用户无法获得我的任务"))
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 获取我参加的项目
// @Id M007
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// #Param name query string false "项目名称"
// @Param status query int false "状态（1进行中2完成不传为全部）"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param type query int false "项目类型"
// @Success 200 object response.ListRes{data=[]Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignedprojects [GET]
func WxGetAssignedProjects(c *gin.Context) {
	var filter AssignedProjectFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	count, list, err := projectService.GetAssignedProject(filter, claims.UserID, claims.PositionID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建项目
// @Id M008
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_info body ProjectNew true "项目信息"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects [POST]
func WxNewProject(c *gin.Context) {
	NewProject(c)
}

// @Summary 根据ID获取项目
// @Id M009
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects/:id [GET]
func WxGetProjectByID(c *gin.Context) {
	GetProjectByID(c)

}

// @Summary 根据ID更新项目
// @Id M010
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Param project_info body ProjectUpdate true "项目信息"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects/:id [PUT]
func WxUpdateProject(c *gin.Context) {
	UpdateProject(c)
}

// @Summary 根据ID删除项目
// @Id M011
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects/:id [DELETE]
func WxDeleteProject(c *gin.Context) {
	DeleteProject(c)
}

// @Summary 新建项目报告
// @Id M012
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_info body ProjectReportNew true "项目报告信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id/reports [POST]
func NewProjectReport(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var report ProjectReportNew
	if err := c.ShouldBindJSON(&report); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	report.User = claims.Username
	report.UserID = claims.UserID
	report.OrganizationID = claims.OrganizationID
	projectService := NewProjectService()
	err := projectService.NewProjectReport(uri.ID, report)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新建项目报告
// @Id M013
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_info body ProjectReportNew true "项目报告信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects/:id/reports [POST]
func WxNewProjectReport(c *gin.Context) {
	NewProjectReport(c)
}

// @Summary 项目报告列表
// @Id M014
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param name query string false "名称"
// @Param status query string false "状态all/active"
// @Success 200 object response.SuccessRes{data=[]ProjectReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id/reports [GET]
func GetProjectReportList(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var filter ProjectReportFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.UserID = claims.UserID
	filter.OrganizationID = claims.OrganizationID
	list, err := projectService.GetProjectReportList(uri.ID, filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 项目报告列表
// @Id M015
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param name query string false "名称"
// @Param status query string false "状态all/active"
// @Success 200 object response.SuccessRes{data=[]ProjectReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects/:id/reports [GET]
func WxGetProjectReportList(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var filter ProjectReportFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.UserID = claims.UserID
	filter.OrganizationID = claims.OrganizationID
	list, err := projectService.GetProjectReportList(uri.ID, filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 根据ID获取项目报告
// @Id M016
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "报告ID"
// @Success 200 object response.SuccessRes{data=ProjectReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projectreports/:id [GET]
func GetProjectReportByID(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	projectService := NewProjectService()
	project, err := projectService.GetProjectReportByID(uri.ID, claims.UserID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, project)

}

// @Summary 根据ID删除项目报告
// @Id M017
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目报告ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projectreports/:id [DELETE]
func DeleteProjectReport(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	projectService := NewProjectService()
	err := projectService.DeleteProjectReport(uri.ID, claims.UserID, claims.Username, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 根据ID获取项目报告
// @Id M018
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "报告ID"
// @Success 200 object response.SuccessRes{data=ProjectReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectreports/:id [GET]
func WxGetProjectReportByID(c *gin.Context) {
	GetProjectReportByID(c)

}

// @Summary 根据ID删除项目报告
// @Id M019
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目报告ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectreports/:id [DELETE]
func WxDeleteProjectReport(c *gin.Context) {
	DeleteProjectReport(c)
}

// @Summary 根据ID更新项目报告
// @Id M020
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "报告ID"
// @Param project_info body ProjectReportNew true "报告信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projectreports/:id [PUT]
func UpdateProjectReport(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var report ProjectReportNew
	if err := c.ShouldBindJSON(&report); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	report.User = claims.Username
	report.OrganizationID = claims.OrganizationID
	report.UserID = claims.UserID
	projectService := NewProjectService()
	err := projectService.UpdateProjectReport(uri.ID, report)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID更新项目报告
// @Id M021
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "报告ID"
// @Param project_info body ProjectReportNew true "报告信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectreports/:id [PUT]
func WxUpdateProjectReport(c *gin.Context) {
	UpdateProjectReport(c)
}

// @Summary 新建项目记录
// @Id M022
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_info body ProjectRecordNew true "项目报告信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id/records [POST]
func NewProjectRecord(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var record ProjectRecordNew
	if err := c.ShouldBindJSON(&record); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	record.User = claims.Username
	record.UserID = claims.UserID
	record.OrganizationID = claims.OrganizationID
	projectService := NewProjectService()
	err := projectService.NewProjectRecord(uri.ID, record)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 项目记录列表
// @Id M023
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.SuccessRes{data=[]ProjectRecordResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id/records [GET]
func GetProjectRecordList(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var filter ProjectRecordFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.UserID = claims.UserID
	filter.OrganizationID = claims.OrganizationID
	count, list, err := projectService.GetProjectRecordList(uri.ID, filter, claims.UserType)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取项目记录
// @Id M024
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "记录ID"
// @Success 200 object response.SuccessRes{data=ProjectRecordResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projectrecords/:id [GET]
func GetProjectRecordByID(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	projectService := NewProjectService()
	project, err := projectService.GetProjectRecordByID(uri.ID, claims.UserID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, project)

}

// @Summary 根据ID删除项目记录
// @Id M025
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目记录ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projectrecords/:id [DELETE]
func DeleteProjectRecord(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	projectService := NewProjectService()
	err := projectService.DeleteProjectRecord(uri.ID, claims.UserID, claims.Username, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 根据ID更新项目记录
// @Id M026
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "记录ID"
// @Param project_info body ProjectRecordNew true "记录信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projectrecords/:id [PUT]
func UpdateProjectRecord(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var record ProjectRecordNew
	if err := c.ShouldBindJSON(&record); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	record.User = claims.Username
	record.OrganizationID = claims.OrganizationID
	record.UserID = claims.UserID
	projectService := NewProjectService()
	err := projectService.UpdateProjectRecord(uri.ID, record, claims.UserType)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新建项目记录
// @Id M027
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_info body ProjectRecordNew true "项目记录信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects/:id/records [POST]
func WxNewProjectRecord(c *gin.Context) {
	NewProjectRecord(c)
}

// @Summary 项目记录列表
// @Id M028
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param name query string false "名称"
// @Param status query string false "状态all/active"
// @Success 200 object response.SuccessRes{data=[]ProjectRecordResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projects/:id/records [GET]
func WxGetProjectRecordList(c *gin.Context) {
	GetProjectRecordList(c)
}

// @Summary 根据ID获取项目记录
// @Id M029
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "记录ID"
// @Success 200 object response.SuccessRes{data=ProjectRecordResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectrecords/:id [GET]
func WxGetProjectRecordByID(c *gin.Context) {
	GetProjectRecordByID(c)

}

// @Summary 根据ID删除项目记录
// @Id M030
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目记录ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectrecords/:id [DELETE]
func WxDeleteProjectRecord(c *gin.Context) {
	DeleteProjectRecord(c)
}

// @Summary 根据ID更新项目记录
// @Id M031
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "记录ID"
// @Param project_info body ProjectRecordNew true "记录信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectrecords/:id [PUT]
func WxUpdateProjectRecord(c *gin.Context) {
	UpdateProjectRecord(c)
}

// @Summary 项目记录列表
// @Id M032
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.SuccessRes{data=[]ProjectRecordResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/projects/:id/records [GET]
func PortalGetProjectRecordList(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var filter ProjectRecordFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	count, list, err := projectService.PortalGetProjectRecordList(uri.ID, filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID已阅项目报告
// @Id M033
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "报告ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectreports/:id/views [POST]
func WxViewProjectReport(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	projectService := NewProjectService()
	err := projectService.ViewProjectReport(uri.ID, claims.OrganizationID, claims.UserID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 项目报告未读列表
// @Id M034
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Success 200 object response.SuccessRes{data=[]ProjectReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/projectreports/unread [GET]
func WxGetUnreadReportList(c *gin.Context) {
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	list, err := projectService.GetProjectReportUnreadList(claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}
