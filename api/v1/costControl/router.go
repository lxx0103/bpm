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
	g.PUT("/payments/:id", UpdatePayment)
	g.GET("/payments", GetPaymentList)
	g.GET("/payments/:id", GetPaymentByID)
	g.DELETE("/payments/:id", DeletePayment)

	g.POST("/paymentRequests/:id/deliverys", NewDelivery)
	g.PUT("/deliverys/:id", UpdateDelivery)
	g.GET("/deliverys", GetDeliveryList)
	g.GET("/deliverys/:id", GetDeliveryByID)
	g.DELETE("/deliverys/:id", DeleteDelivery)

	g.POST("/incomes", NewIncome)
	g.PUT("/incomes/:id", UpdateIncome)
	g.GET("/incomes", GetIncomeList)
	g.GET("/incomes/:id", GetIncomeByID)
	g.DELETE("/incomes/:id", DeleteIncome)

	g.GET("/reports/project/:id", GetReportByProjectID)
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
	g.GET("/wx/paymentRequestHistorys", WxGetPaymentRequestHistory)
	g.PUT("/wx/paymentRequests/:id/audit", WxUpdatePaymentRequestAudit)

	g.POST("/wx/paymentRequests/:id/payments", WxNewPayment)
	g.PUT("/wx/payments/:id", WxUpdatePayment)
	g.GET("/wx/payments", WxGetPaymentList)
	g.GET("/wx/payments/:id", WxGetPaymentByID)
	g.DELETE("/wx/payments/:id", WxDeletePayment)

	g.POST("/wx/incomes", WxNewIncome)
	g.PUT("/wx/incomes/:id", WxUpdateIncome)
	g.GET("/wx/incomes", WxGetIncomeList)
	g.GET("/wx/incomes/:id", WxGetIncomeByID)
	g.DELETE("/wx/incomes/:id", WxDeleteIncome)

	g.POST("/wx/paymentRequests/:id/deliverys", WxNewDelivery)
	g.PUT("/wx/deliverys/:id", WxUpdateDelivery)
	g.GET("/wx/deliverys", WxGetDeliveryList)
	g.GET("/wx/deliverys/:id", WxGetDeliveryByID)
	g.DELETE("/wx/deliverys/:id", WxDeleteDelivery)

	g.GET("/wx/reports/project/:id", WxGetReportByProjectID)
}
