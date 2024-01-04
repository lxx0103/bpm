package costControl

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 项目预算清单
// @Id S001
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int true "项目ID"
// @Param organization_id query int false "组织名称"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.ListRes{data=[]RespBudget} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /budgets [GET]
func GetBudgetList(c *gin.Context) {
	var filter ReqBudgetFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	err = filter.Verify()
	if err != nil {
		response.ResponseError(c, "VerifyError", err)
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	if claims.OrganizationID != 0 {
		filter.OrganizationID = claims.OrganizationID
	}
	costControlService := NewCostControlService()
	count, list, err := costControlService.GetBudgetList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建项目预算
// @Id S002
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body ReqBudgetNew true "预算信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /budgets [POST]
func NewBudget(c *gin.Context) {
	var info ReqBudgetNew
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	err = info.Verify()
	if err != nil {
		response.ResponseError(c, "VerifyError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	if claims.OrganizationID != 0 {
		info.OrganizationID = claims.OrganizationID
	}
	info.User = claims.Username
	info.UserID = claims.UserID
	costControlService := NewCostControlService()
	err = costControlService.NewBudget(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 更新项目预算
// @Id S003
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "预算ID"
// @Param info body ReqBudgetUpdate true "预算信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /budgets/:id [PUT]
func UpdateBudget(c *gin.Context) {
	var uri BudgetID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqBudgetUpdate
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	err = info.Verify()
	if err != nil {
		response.ResponseError(c, "VerifyError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	costControlService := NewCostControlService()
	err = costControlService.UpdateBudget(info, uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取项目预算
// @Id S004
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "预算ID"
// @Success 200 object response.SuccessRes{data=RespBudget} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /budgets/:id [GET]
func GetBudgetByID(c *gin.Context) {
	var uri BudgetID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	row, err := costControlService.GetBudgetByID(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, row)
}

// @Summary 删除项目预算
// @Id S005
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "预算ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /budgets/:id [DELETE]
func DeleteBudget(c *gin.Context) {
	var uri BudgetID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	err := costControlService.DeleteBudget(uri.ID, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新建费用申请
// @Id S006
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body ReqPaymentRequestNew true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests [POST]
func NewPaymentRequest(c *gin.Context) {
	var info ReqPaymentRequestNew
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	err = info.Verify()
	if err != nil {
		response.ResponseError(c, "VerifyError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	if claims.OrganizationID != 0 {
		info.OrganizationID = claims.OrganizationID
	}
	info.User = claims.Username
	info.UserID = claims.UserID
	costControlService := NewCostControlService()
	err = costControlService.NewPaymentRequest(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 更新费用申请
// @Id S007
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentRequestUpdate true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests/:id [PUT]
func UpdatePaymentRequest(c *gin.Context) {
	var uri PaymentRequestID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqPaymentRequestUpdate
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	err = info.Verify()
	if err != nil {
		response.ResponseError(c, "VerifyError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	costControlService := NewCostControlService()
	err = costControlService.UpdatePaymentRequest(info, uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 费用申请列表
// @Id S008
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int false "项目ID"
// @Param organization_id query int false "组织ID"
// @Param name query string false "名称"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.ListRes{data=[]RespPaymentRequest} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests [GET]
func GetPaymentRequestList(c *gin.Context) {
	var filter ReqPaymentRequestFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	err = filter.Verify()
	if err != nil {
		response.ResponseError(c, "VerifyError", err)
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	if claims.OrganizationID != 0 {
		filter.OrganizationID = claims.OrganizationID
	}
	costControlService := NewCostControlService()
	count, list, err := costControlService.GetPaymentRequestList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取费用申请
// @Id S009
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Success 200 object response.SuccessRes{data=RespPaymentRequest} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests/:id [GET]
func GetPaymentRequestByID(c *gin.Context) {
	var uri PaymentRequestID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	row, err := costControlService.GetPaymentRequestByID(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, row)
}

// @Summary 删除费用申请
// @Id S010
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests/:id [DELETE]
func DeletePaymentRequest(c *gin.Context) {
	var uri PaymentRequestID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	err := costControlService.DeletePaymentRequest(uri.ID, claims.OrganizationID, claims.Username, claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 更新费用申请审核设置
// @Id S011
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body ReqPaymentRequestTypeUpdate true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequestTypes [POST]
func UpdatePaymentRequestType(c *gin.Context) {
	var info ReqPaymentRequestTypeUpdate
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	err = info.Verify()
	if err != nil {
		response.ResponseError(c, "VerifyError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	costControlService := NewCostControlService()
	err = costControlService.UpdatePaymentRequestType(info, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 费用申请审核设置列表
// @Id S012
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param organization_id query int false "组织ID"
// @Success 200 object response.SuccessRes{data=[]RespPaymentRequestType} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequestTypes [GET]
func GetPaymentRequestTypeList(c *gin.Context) {
	var filter ReqPaymentRequestTypeFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	if claims.OrganizationID != 0 {
		filter.OrganizationID = claims.OrganizationID
	}
	costControlService := NewCostControlService()
	res, err := costControlService.GetPaymentRequestTypeList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}
