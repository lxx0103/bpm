package costControl

import "errors"

type ReqBudgetNew struct {
	OrganizationID int64    `json:"organization_id" binding:"required,min=1"`
	ProjectID      int64    `json:"project_id" binding:"required,min=1"`
	BudgetType     int      `json:"budget_type" binding:"required,min=1,max=2"`
	Name           string   `json:"name" binding:"required,max=100"`
	Quantity       int64    `json:"quantity" binding:"required,min=1"`
	UnitPrice      float64  `json:"unit_price" binding:"required,min=0"`
	Budget         float64  `json:"budget" binding:"required,min=0"`
	Remark         string   `json:"remark" binding:"omitempty,max=255"`
	Picture        []string `json:"picture" binding:"omitempty"`
	User           string   `json:"user" swaggerignore:"true"`
	UserID         int64    `json:"user_id" swaggerignore:"true"`
}

func (f *ReqBudgetNew) Verify() error {
	if f.Budget != float64(f.Quantity)*f.UnitPrice {
		msg := "总预算错误"
		return errors.New(msg)
	}
	return nil
}

type ReqBudgetPictureNew struct {
	BudgetID int64  `json:"budget_id" binding:"required,min=1"`
	Picture  string `json:"picture" binding:"required"`
	User     string `json:"user" swaggerignore:"true"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type BudgetID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ReqBudgetUpdate struct {
	BudgetType int      `json:"budget_type" binding:"required,min=1,max=2"`
	Name       string   `json:"name" binding:"required,max=100"`
	Quantity   int64    `json:"quantity" binding:"required,min=1"`
	UnitPrice  float64  `json:"unit_price" binding:"required,min=0"`
	Budget     float64  `json:"budget" binding:"required,min=0"`
	Used       float64  `json:"used" swaggerignore:"true"`
	Balance    float64  `json:"balance" swaggerignore:"true"`
	Remark     string   `json:"remark" binding:"omitempty,max=255"`
	Picture    []string `json:"picture" binding:"omitempty"`
	User       string   `json:"user" swaggerignore:"true"`
	UserID     int64    `json:"user_id" swaggerignore:"true"`
}

func (f *ReqBudgetUpdate) Verify() error {
	if f.Budget != float64(f.Quantity)*f.UnitPrice {
		msg := "总预算错误"
		return errors.New(msg)
	}
	return nil
}

type ReqBudgetFilter struct {
	ProjectID      int64  `form:"project_id" binding:"required,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	BudgetType     int    `form:"budget_type" binding:"omitempty,min=1,max=2"`
	Name           string `form:"name" binding:"omitempty,max=100"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

func (f *ReqBudgetFilter) Verify() error {
	return nil
}

type RespBudget struct {
	ID               int64    `db:"id" json:"id"`
	OrganizationID   int64    `db:"organization_id" json:"organization_id"`
	OrganizationName string   `db:"organization_name" json:"organization_name"`
	ProjectID        int64    `db:"project_id" json:"project_id"`
	ProjectName      string   `db:"project_name" json:"project_name"`
	BudgetType       int      `db:"budget_type" json:"budget_type"`
	Name             string   `db:"name" json:"name"`
	Quantity         int64    `db:"quantity" json:"quantity"`
	UnitPrice        float64  `db:"unit_price" json:"unit_price"`
	Budget           float64  `db:"budget" json:"budget"`
	Used             float64  `db:"used" json:"used"`
	Balance          float64  `db:"balance" json:"balance"`
	Remark           string   `db:"remark" json:"remark"`
	Picture          []string `json:"picture"`
	Status           int      `db:"status" json:"status"`
}

type ReqPaymentRequestNew struct {
	OrganizationID     int64    `json:"organization_id" binding:"required,min=1"`
	ProjectID          int64    `json:"project_id" binding:"omitempty"`
	PaymentRequestType int      `json:"payment_request_type" binding:"required,min=1,max=2"` // 1:采购类，2:工款类
	BudgetID           int64    `json:"budget_id" binding:"omitempty,min=1"`                 //可选Budget（如果project id 不为空）
	Name               string   `json:"name" binding:"required,max=100"`
	Quantity           int64    `json:"quantity" binding:"required,min=1"`
	UnitPrice          float64  `json:"unit_price" binding:"required"`
	Total              float64  `json:"total" binding:"required"`
	Remark             string   `json:"remark" binding:"omitempty,max=255"`
	Picture            []string `json:"picture" binding:"omitempty"`
	User               string   `json:"user" swaggerignore:"true"`
	UserID             int64    `json:"user_id" swaggerignore:"true"`
}

func (f *ReqPaymentRequestNew) Verify() error {
	if f.Total != float64(f.Quantity)*f.UnitPrice {
		msg := "总费用错误"
		return errors.New(msg)
	}
	return nil
}

type PaymentRequestID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ReqPaymentRequestUpdate struct {
	ProjectID          int64    `json:"project_id" binding:"omitempty"`
	PaymentRequestType int      `json:"payment_request_type" binding:"required,min=1,max=2"` // 1:采购类，2:工款类
	BudgetID           int64    `json:"budget_id" binding:"omitempty,min=1"`                 //可选Budget（如果project id 不为空）
	Name               string   `json:"name" binding:"required,max=100"`
	Quantity           int64    `json:"quantity" binding:"required,min=1"`
	UnitPrice          float64  `json:"unit_price" binding:"required"`
	Total              float64  `json:"total" binding:"required"`
	Remark             string   `json:"remark" binding:"omitempty,max=255"`
	Picture            []string `json:"picture" binding:"omitempty"`
	Status             int      `json:"status" swaggerignore:"true"`
	User               string   `json:"user" swaggerignore:"true"`
	UserID             int64    `json:"user_id" swaggerignore:"true"`
}

func (f *ReqPaymentRequestUpdate) Verify() error {
	if f.Total != float64(f.Quantity)*f.UnitPrice {
		msg := "总费用错误"
		return errors.New(msg)
	}
	return nil
}

type ReqPaymentRequestHistoryNew struct {
	PaymentRequestID int64  `json:"payment_request_id" binding:"required,min=1"`
	OrganizationID   int64  `json:"organization_id" binding:"required,min=1"`
	Action           string `json:"action" binding:"required"`
	Remark           string `json:"remark" binding:"omitempty"`
	User             string `json:"user" swaggerignore:"true"`
	UserID           int64  `json:"user_id" swaggerignore:"true"`
}

type ReqPaymentRequestPictureNew struct {
	PaymentRequestID int64  `json:"payment_request_id" binding:"required,min=1"`
	Picture          string `json:"picture" binding:"required"`
	User             string `json:"user" swaggerignore:"true"`
	UserID           int64  `json:"user_id" swaggerignore:"true"`
}

type RespPaymentRequest struct {
	ID                 int64    `db:"id" json:"id"`
	OrganizationID     int64    `db:"organization_id" json:"organization_id"`
	OrganizationName   string   `db:"organization_name" json:"organization_name"`
	ProjectID          int64    `db:"project_id" json:"project_id"`
	ProjectName        string   `db:"project_name" json:"project_name"`
	PaymentRequestType int      `db:"payment_request_type" json:"payment_request_type"`
	BudgetID           int64    `db:"budget_id" json:"budget_id"`
	Name               string   `db:"name" json:"name"`
	Quantity           int64    `db:"quantity" json:"quantity"`
	UnitPrice          float64  `db:"unit_price" json:"unit_price"`
	Total              float64  `db:"total" json:"total"`
	Paid               float64  `db:"paid" json:"paid"`
	Due                float64  `db:"due" json:"due"`
	Remark             string   `db:"remark" json:"remark"`
	Picture            []string `json:"picture"`
	UserID             int64    `db:"user_id" json:"user_id"`
	Status             int      `db:"status" json:"status"`
}

type ReqPaymentRequestFilter struct {
	ProjectID          int64  `form:"project_id" binding:"omitempty,min=1"`
	OrganizationID     int64  `form:"organization_id" binding:"omitempty,min=1"`
	PaymentRequestType int    `form:"payment_request_type" binding:"omitempty,min=1,max=2"`
	Name               string `form:"name" binding:"omitempty,max=100"`
	PageId             int    `form:"page_id" binding:"required,min=1"`
	PageSize           int    `form:"page_size" binding:"required,min=5,max=200"`
}

func (f *ReqPaymentRequestFilter) Verify() error {
	return nil
}

type ReqPaymentRequestTypeUpdate struct {
	ReqPaymentRequestType int `json:"payment_request_type" binding:"required,min=1,max=2"`
	AuditInfo             []struct {
		AuditLevel int     `json:"audit_level" binding:"required,min=1"`
		AuditType  int     `json:"audit_type" binding:"required,oneof=1 2"`
		AuditTo    []int64 `json:"audit_to" binding:"required"`
	} `json:"audit_info" binding:"omitempty"`
	User   string `json:"user" swaggerignore:"true"`
	UserID int64  `json:"user_id" swaggerignore:"true"`
}

func (f *ReqPaymentRequestTypeUpdate) Verify() error {
	if len(f.AuditInfo) == 0 {
		msg := "必须至少有一层审核"
		return errors.New(msg)
	}
	return nil
}

type ReqPaymentRequestTypeAudit struct {
	OrganizationID     int64  `json:"organization_id"`
	PaymentRequestType int    `json:"payment_request_type"`
	AuditLevel         int    `json:"audit_level"`
	AuditType          int    `json:"audit_type"`
	AuditTo            int64  `json:"audit_to"`
	User               string `json:"user" swaggerignore:"true"`
	UserID             int64  `json:"user_id" swaggerignore:"true"`
}

type ReqPaymentRequestTypeFilter struct {
	OrganizationID int64 `form:"organization_id" binding:"omitempty,min=1"`
}

type RespPaymentRequestType struct {
	PaymentRequestType     int                           `json:"payment_request_type"`
	PaymentRequestTypeName string                        `json:"payment_request_type_name"`
	Audit                  []RespPaymentRequestTypeAudit `json:"audit"`
}

type RespPaymentRequestTypeAudit struct {
	AuditLevel  int    `db:"audit_level" json:"audit_level"`
	AuditType   int    `db:"audit_type" json:"audit_type"`
	AuditTo     int64  `db:"audit_to" json:"audit_to"`
	AuditToName string `db:"audit_to_name" json:"audit_to_name"`
}
