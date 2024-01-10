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
