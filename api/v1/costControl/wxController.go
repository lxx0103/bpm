package costControl

import "github.com/gin-gonic/gin"

// @Summary 项目预算清单
// @Id WXS001
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int true "项目ID"
// @Param organization_id query int false "组织名称"
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Success 200 object response.ListRes{data=[]RespBudget} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/budgets [GET]
func WxGetBudgetList(c *gin.Context) {
	GetBudgetList(c)
}

// @Summary 新建项目预算
// @Id WXS002
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body ReqBudgetNew true "预算信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/budgets [POST]
func WxNewBudget(c *gin.Context) {
	NewBudget(c)
}

// @Summary 更新项目预算
// @Id WXS003
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "预算ID"
// @Param info body ReqBudgetUpdate true "预算信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/budgets/:id [PUT]
func WxUpdateBudget(c *gin.Context) {
	UpdateBudget(c)
}

// @Summary 根据ID获取项目预算
// @Id WXS004
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "预算ID"
// @Success 200 object response.SuccessRes{data=RespBudget} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/budgets/:id [GET]
func WxGetBudgetByID(c *gin.Context) {
	GetBudgetByID(c)
}

// @Summary 删除项目预算
// @Id WXS005
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "预算ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/budgets/:id [DELETE]
func WxDeleteBudget(c *gin.Context) {
	DeleteBudget(c)
}

// @Summary 新建费用申请
// @Id WXS006
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body ReqPaymentRequestNew true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequests [POST]
func WxNewPaymentRequest(c *gin.Context) {
	NewPaymentRequest(c)
}

// @Summary 更新费用申请
// @Id WXS007
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentRequestUpdate true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequests/:id [PUT]
func WxUpdatePaymentRequest(c *gin.Context) {
	UpdatePaymentRequest(c)
}

// @Summary 费用申请列表
// @Id WXS008
// @Tags 小程序成控管理
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
// @Router /wx/paymentRequests [GET]
func WxGetPaymentRequestList(c *gin.Context) {
	GetPaymentRequestList(c)
}

// @Summary 根据ID获取费用申请
// @Id WXS009
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Success 200 object response.SuccessRes{data=RespPaymentRequest} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequests/:id [GET]
func WxGetPaymentRequestByID(c *gin.Context) {
	GetPaymentRequestByID(c)
}

// @Summary 删除费用申请
// @Id WXS010
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequests/:id [DELETE]
func WxDeletePaymentRequest(c *gin.Context) {
	DeletePaymentRequest(c)
}

// @Summary 更新费用申请审核设置
// @Id WXS011
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body ReqPaymentRequestTypeUpdate true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequestTypes [POST]
func WxUpdatePaymentRequestType(c *gin.Context) {
	UpdatePaymentRequestType(c)
}

// @Summary 费用申请审核设置列表
// @Id WXS012
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param organization_id query int false "组织ID"
// @Success 200 object response.SuccessRes{data=[]RespPaymentRequestType} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequestTypes [GET]
func WxGetPaymentRequestTypeList(c *gin.Context) {
	GetPaymentRequestTypeList(c)
}

// @Summary 审核费用申请
// @Id WXS013
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentRequestAudit true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequests/:id/audit [POST]
func WxAuditPaymentRequest(c *gin.Context) {
	AuditPaymentRequest(c)
}

// @Summary 费用申请操作历史列表
// @Id WXS014
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param payment_request_id query int true "请款ID"
// @Param organization_id query int false "组织ID"
// @Success 200 object response.ListRes{data=[]RespPaymentRequest} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequestHistorys [GET]
func WxGetPaymentRequestHistory(c *gin.Context) {
	GetPaymentRequestHistory(c)
}

// @Summary 更新费用申请审核
// @Id WXS015
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentRequestAuditUpdate true "请款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequests/:id/audit [PUT]
func WxUpdatePaymentRequestAudit(c *gin.Context) {
	UpdatePaymentRequestAudit(c)
}

// @Summary 新增付款
// @Id WXS016
// @Tags 小程序成控管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "请款ID"
// @Param info body ReqPaymentNew true "付款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/paymentRequests/:id/payments [POST]
func WxNewPayment(c *gin.Context) {
	NewPayment(c)
}