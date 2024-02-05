package costControl

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/budgets", NewBudget)
	g.PUT("/budgets/:id", UpdateBudget)
	g.GET("/budgets", GetBudgetList)
	g.GET("/budgets/:id", GetBudgetByID)
	g.DELETE("/budgets/:id", DeleteBudget)

	g.POST("/paymentRequests", NewPaymentRequest)
	g.PUT("/paymentRequests/:id", UpdatePaymentRequest)
	g.GET("/paymentRequests", GetPaymentRequestList)
	g.GET("/paymentRequests/:id", GetPaymentRequestByID)
	g.DELETE("/paymentRequests/:id", DeletePaymentRequest)
	g.POST("/paymentRequests/:id/audit", AuditPaymentRequest)
	g.GET("paymentRequestHistorys", GetPaymentRequestHistory)
	g.PUT("/paymentRequests/:id/audit", UpdatePaymentRequestAudit)

	g.POST("/paymentRequestTypes", UpdatePaymentRequestType)
	g.GET("/paymentRequestTypes", GetPaymentRequestTypeList)

	g.POST("/paymentRequests/:id/payments", NewPayment)
	// g.PUT("/payments/:id", UpdatePayment)
	// g.GET("/payments", GetPaymentList)
	// g.GET("/payments/:id", GetPaymentByID)
	// g.DELETE("/payments/:id", DeletePayment)

}

func WxRouters(g *gin.RouterGroup) {
	g.POST("/wx/budgets", WxNewBudget)
	g.PUT("/wx/budgets/:id", WxUpdateBudget)
	g.GET("/wx/budgets", WxGetBudgetList)
	g.GET("/wx/budgets/:id", WxGetBudgetByID)
	g.DELETE("/wx/budgets/:id", WxDeleteBudget)

	g.POST("/wx/paymentRequests", WxNewPaymentRequest)
	g.PUT("/wx/paymentRequests/:id", WxUpdatePaymentRequest)
	g.GET("/wx/paymentRequests", WxGetPaymentRequestList)
	g.GET("/wx/paymentRequests/:id", WxGetPaymentRequestByID)
	g.DELETE("/wx/paymentRequests/:id", WxDeletePaymentRequest)
	g.POST("/wx/paymentRequests/:id/audit", WxAuditPaymentRequest)
	g.GET("wxpaymentRequestHistorys", WxGetPaymentRequestHistory)
	g.PUT("/wx/paymentRequests/:id/audit", WxUpdatePaymentRequestAudit)

	g.POST("/wx/paymentRequests/:id/payments", WxNewPayment)
}