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

	g.POST("/paymentRequestTypes", UpdatePaymentRequestType)
	g.GET("/paymentRequestTypes", GetPaymentRequestTypeList)

}
