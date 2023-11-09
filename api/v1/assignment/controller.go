package assignment

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 任务列表
// @Id Q001
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "任务名称"
// @Param organization_id query int64 false "组织ID"
// @Param assignment_type query int64 false "任务类型（1，会议任务，2其他任务）"
// @Param reference_id query int64 false "关联ID（会议ID)"
// @Param project_id query int64 false "项目ID"
// @Param event_id query int64 false "事件ID"
// @Success 200 object response.ListRes{data=[]AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments [GET]
func GetAssignmentList(c *gin.Context) {
	var filter AssignmentFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	assignmentService := NewAssignmentService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := assignmentService.GetAssignmentList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建任务
// @Id Q002
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param assignment_info body AssignmentNew true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments [POST]
func NewAssignment(c *gin.Context) {
	var assignment AssignmentNew
	if err := c.ShouldBindJSON(&assignment); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	assignment.User = claims.Username
	assignment.UserID = claims.UserID
	organizationID := claims.OrganizationID
	assignmentService := NewAssignmentService()
	err := assignmentService.NewAssignment(assignment, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取任务
// @Id Q003
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Success 200 object response.SuccessRes{data=AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/:id [GET]
func GetAssignmentByID(c *gin.Context) {
	var uri AssignmentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	assignmentService := NewAssignmentService()
	assignment, err := assignmentService.GetAssignmentByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, assignment)

}

// @Summary 根据ID更新任务
// @Id Q004
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Param info body AssignmentUpdate true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/:id [PUT]
func UpdateAssignment(c *gin.Context) {
	var uri AssignmentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var assignment AssignmentUpdate
	if err := c.ShouldBindJSON(&assignment); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	assignment.User = claims.Username
	assignment.UserID = claims.UserID
	organizationID := claims.OrganizationID
	assignmentService := NewAssignmentService()
	err := assignmentService.UpdateAssignment(uri.ID, assignment, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除任务
// @Id Q005
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/:id [DELETE]
func DeleteAssignment(c *gin.Context) {
	var uri AssignmentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	assignmentService := NewAssignmentService()
	err := assignmentService.DeleteAssignment(uri.ID, claims.OrganizationID, claims.Username, claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID完成任务
// @Id Q006
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Param info body AssignmentComplete true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/:id/complete [POST]
func CompleteAssignment(c *gin.Context) {
	var uri AssignmentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info AssignmentComplete
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	organizationID := claims.OrganizationID
	assignmentService := NewAssignmentService()
	err := assignmentService.CompleteAssignment(uri.ID, info, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID审核任务
// @Id Q007
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Param info body AssignmentAudit true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/:id/audit [POST]
func AuditAssignment(c *gin.Context) {
	var uri AssignmentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info AssignmentAudit
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	organizationID := claims.OrganizationID
	assignmentService := NewAssignmentService()
	err := assignmentService.AuditAssignment(uri.ID, info, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 我的任务列表
// @Id Q008
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "任务名称"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/my [GET]
func GetMyAssignmentList(c *gin.Context) {
	var filter MyAssignmentFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	assignmentService := NewAssignmentService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.UserID = claims.UserID
	count, list, err := assignmentService.GetMyAssignmentList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 我的任务审核列表
// @Id Q009
// @Tags 任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "任务名称"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/myaudit [GET]
func GetMyAuditList(c *gin.Context) {
	var filter MyAuditFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	assignmentService := NewAssignmentService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.UserID = claims.UserID
	count, list, err := assignmentService.GetMyAuditList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 任务列表
// @Id Q010
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "任务名称"
// @Param organization_id query int64 false "组织ID"
// @Param assignment_type query int64 false "任务类型（1，会议任务，2其他任务）"
// @Param reference_id query int64 false "关联ID（会议ID)"
// @Param project_id query int64 false "项目ID"
// @Success 200 object response.ListRes{data=[]AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments [GET]
func WxGetAssignmentList(c *gin.Context) {
	GetAssignmentList(c)
}

// @Summary 新建任务
// @Id Q011
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param assignment_info body AssignmentNew true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments [POST]
func WxNewAssignment(c *gin.Context) {
	NewAssignment(c)
}

// @Summary 根据ID获取任务
// @Id Q012
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Success 200 object response.SuccessRes{data=AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments/:id [GET]
func WxGetAssignmentByID(c *gin.Context) {
	GetAssignmentByID(c)
}

// @Summary 根据ID更新任务
// @Id Q013
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Param info body AssignmentUpdate true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments/:id [PUT]
func WxUpdateAssignment(c *gin.Context) {
	UpdateAssignment(c)
}

// @Summary 根据ID删除任务
// @Id Q014
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments/:id [DELETE]
func WxDeleteAssignment(c *gin.Context) {
	DeleteAssignment(c)
}

// @Summary 根据ID完成任务
// @Id Q015
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Param info body AssignmentComplete true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments/:id/complete [POST]
func WxCompleteAssignment(c *gin.Context) {
	CompleteAssignment(c)
}

// @Summary 根据ID审核任务
// @Id Q016
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "任务ID"
// @Param info body AssignmentAudit true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments/:id/audit [POST]
func WxAuditAssignment(c *gin.Context) {
	AuditAssignment(c)
}

// @Summary 我的任务列表
// @Id Q017
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "任务名称"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments/my [GET]
func WxGetMyAssignmentList(c *gin.Context) {
	GetMyAssignmentList(c)
}

// @Summary 我的任务审核列表
// @Id Q018
// @Tags 小程序接口-任务管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "任务名称"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]AssignmentResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments/myaudit [GET]
func WxGetMyAuditList(c *gin.Context) {
	GetMyAuditList(c)
}
