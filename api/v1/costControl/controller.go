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
// @Param type query string false "类型（audit：审核人员， mine：我创建的，passed：已审核通过的）"
// @Param payment_status query string false "付款状态（none：未付款， partial：部分付款，paid：已付款）"
// @Param delivery_status query string false "进场状态（none：未进场， partial：部分进场，deliveried：已进场）"
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
	filter.User = claims.Username
	filter.UserID = claims.UserID
	filter.PositionID = claims.PositionID
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
	if claims.OrganizationID != 0 {
		info.OrganizationID = claims.OrganizationID
	}
	costControlService := NewCostControlService()
	err = costControlService.UpdatePaymentRequestType(info)
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

// @Summary 审核费用申请
// @Id S013
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentRequestAudit true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests/:id/audit [POST]
func AuditPaymentRequest(c *gin.Context) {
	var uri PaymentRequestID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqPaymentRequestAudit
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.PositionID = claims.PositionID
	costControlService := NewCostControlService()
	err = costControlService.AuditPaymentRequest(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 费用申请操作历史列表
// @Id S014
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param payment_request_id query int true "请款ID"
// @Param organization_id query int false "组织ID"
// @Success 200 object response.ListRes{data=[]RespPaymentRequest} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequestHistorys [GET]
func GetPaymentRequestHistory(c *gin.Context) {
	var filter ReqPaymentRequestHistoryFilter
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
	list, err := costControlService.GetPaymentRequestHistoryList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 更新费用申请审核
// @Id S015
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentRequestAuditUpdate true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests/:id/audit [PUT]
func UpdatePaymentRequestAudit(c *gin.Context) {
	var uri PaymentRequestID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqPaymentRequestAuditUpdate
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
	err = costControlService.UpdatePaymentRequestAudit(uri.ID, info, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新增付款
// @Id S016
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentNew true "付款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests/:id/payments [POST]
func NewPayment(c *gin.Context) {
	var uri PaymentRequestID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqPaymentNew
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	costControlService := NewCostControlService()
	err = costControlService.NewPayment(uri.ID, info, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 更新付款
// @Id S017
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "付款ID"
// @Param info body ReqPaymentUpdate true "付款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /payments/:id [PUT]
func UpdatePayment(c *gin.Context) {
	var uri PaymentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqPaymentUpdate
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.OrganizationID = claims.OrganizationID
	costControlService := NewCostControlService()
	err = costControlService.UpdatePayment(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 付款列表
// @Id S018
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int false "项目ID"
// @Param payment_request_id query int false "请款ID"
// @Param organization_id query int false "组织ID"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.ListRes{data=[]RespPayment} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /payments [GET]
func GetPaymentList(c *gin.Context) {
	var filter ReqPaymentFilter
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
	count, list, err := costControlService.GetPaymentList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取付款
// @Id S019
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "付款ID"
// @Success 200 object response.SuccessRes{data=RespPayment} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /payments/:id [GET]
func GetPaymentByID(c *gin.Context) {
	var uri PaymentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	row, err := costControlService.GetPaymentByID(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, row)
}

// @Summary 删除付款
// @Id S020
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "付款ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /payments/:id [DELETE]
func DeletePayment(c *gin.Context) {
	var uri PaymentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	err := costControlService.DeletePayment(uri.ID, claims.OrganizationID, claims.Username, claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新增收入
// @Id S021
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body ReqIncomeNew true "收入信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /incomes [POST]
func NewIncome(c *gin.Context) {
	var info ReqIncomeNew
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.OrganizationID = claims.OrganizationID
	costControlService := NewCostControlService()
	err = costControlService.NewIncome(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 更新收入
// @Id S022
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "收入ID"
// @Param info body ReqIncomeUpdate true "收入信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /incomes/:id [PUT]
func UpdateIncome(c *gin.Context) {
	var uri IncomeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqIncomeUpdate
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.OrganizationID = claims.OrganizationID
	costControlService := NewCostControlService()
	err = costControlService.UpdateIncome(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 收入列表
// @Id S023
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int false "项目ID"
// @Param organization_id query int false "组织ID"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.ListRes{data=[]RespIncome} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /incomes [GET]
func GetIncomeList(c *gin.Context) {
	var filter ReqIncomeFilter
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
	count, list, err := costControlService.GetIncomeList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取收入
// @Id S024
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "收入ID"
// @Success 200 object response.SuccessRes{data=RespIncome} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /incomes/:id [GET]
func GetIncomeByID(c *gin.Context) {
	var uri IncomeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	row, err := costControlService.GetIncomeByID(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, row)
}

// @Summary 删除收入
// @Id S025
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "收入ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /incomes/:id [DELETE]
func DeleteIncome(c *gin.Context) {
	var uri IncomeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	err := costControlService.DeleteIncome(uri.ID, claims.OrganizationID, claims.Username, claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新增材料进场
// @Id S026
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqDeliveryNew true "材料进场信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentRequests/:id/deliverys [POST]
func NewDelivery(c *gin.Context) {
	var uri PaymentRequestID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqDeliveryNew
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	costControlService := NewCostControlService()
	err = costControlService.NewDelivery(uri.ID, info, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 更新材料进场
// @Id S027
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "材料进场ID"
// @Param info body ReqDeliveryUpdate true "材料进场信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /deliverys/:id [PUT]
func UpdateDelivery(c *gin.Context) {
	var uri DeliveryID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ReqDeliveryUpdate
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.OrganizationID = claims.OrganizationID
	costControlService := NewCostControlService()
	err = costControlService.UpdateDelivery(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 材料进场列表
// @Id S028
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int false "项目ID"
// @Param organization_id query int false "组织ID"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.ListRes{data=[]RespDelivery} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /deliverys [GET]
func GetDeliveryList(c *gin.Context) {
	var filter ReqDeliveryFilter
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
	count, list, err := costControlService.GetDeliveryList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取材料进场
// @Id S029
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "材料进场ID"
// @Success 200 object response.SuccessRes{data=RespDelivery} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /deliverys/:id [GET]
func GetDeliveryByID(c *gin.Context) {
	var uri DeliveryID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	row, err := costControlService.GetDeliveryByID(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, row)
}

// @Summary 删除材料进场
// @Id S030
// @Tags 成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "材料进场ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /deliverys/:id [DELETE]
func DeleteDelivery(c *gin.Context) {
	var uri DeliveryID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	costControlService := NewCostControlService()
	err := costControlService.DeleteDelivery(uri.ID, claims.OrganizationID, claims.Username, claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}
