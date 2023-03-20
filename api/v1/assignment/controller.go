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
// @Param assignment_info body AssignmentNew true "任务信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /assignments/:id [PUT]
func UpdateAssignment(c *gin.Context) {
	var uri AssignmentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var assignment AssignmentNew
	if err := c.ShouldBindJSON(&assignment); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	assignment.User = claims.Username
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
	err := assignmentService.DeleteAssignment(uri.ID, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 任务列表
// @Id Q006
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "任务编码"
// @Success 200 object response.ListRes{data=[]Assignment} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/assignments [GET]
func WxGetAssignmentList(c *gin.Context) {
	GetAssignmentList(c)
}
