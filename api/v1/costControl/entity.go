package costControl

import "time"

type Budget struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	ProjectID      int64     `db:"project_id" json:"project_id"`
	BudgetType     int       `db:"budget_type" json:"budget_type"` //1. 材料 2. 人工
	Name           string    `db:"name" json:"name"`
	Quantity       int       `db:"quantity" json:"quantity"`
	UnitPrice      float64   `db:"unit_price" json:"unit_price"`
	Budget         float64   `db:"budget" json:"budget"`
	Used           float64   `db:"used" json:"used"`
	Balance        float64   `db:"balance" json:"balance"`
	Remark         string    `db:"remark" json:"remark"`
	Picture        string    `db:"picture" json:"picture"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type BudgetPicture struct {
	ID        int64     `db:"id" json:"id"`
	BudgetID  int64     `db:"budget_id" json:"budget_id"`
	Link      string    `db:"link" json:"link"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type BudgetDetail struct {
	ID               int64     `db:"id" json:"id"`
	BudgetID         int64     `db:"budget_id" json:"budget_id"`
	PaymentRequestID int64     `db:"payment_request_id"`
	Amount           float64   `db:"amount" json:"amount"`
	Status           int       `db:"status" json:"status"`
	Created          time.Time `db:"created" json:"created"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	Updated          time.Time `db:"updated" json:"updated"`
	UpdatedBy        string    `db:"updated_by" json:"updated_by"`
}

type PaymentRequest struct {
	ID                 int64     `db:"id" json:"id"`
	OrganizationID     int64     `db:"organization_id" json:"organization_id"`
	ProjectID          int64     `db:"project_id" json:"project_id"`
	PaymentRequestType int       `db:"payment_request_type" json:"payment_request_type"` //1. 采购 2. 工款
	BudgetID           int64     `db:"budget_id" json:"budget_id"`
	Name               string    `db:"name" json:"name"`
	Quantity           int       `db:"quantity" json:"quantity"`
	UnitPrice          float64   `db:"unit_price" json:"unit_price"`
	Total              float64   `db:"total" json:"total"`
	Paid               float64   `db:"paid" json:"paid"`
	Due                float64   `db:"due" json:"due"`
	Remark             string    `db:"remark" json:"remark"`
	Picture            string    `db:"picture" json:"picture"`
	Status             int       `db:"status" json:"status"` // 1.待审核，2.审核通过，3.审核驳回，4.部分付款，5.已付款，-1.删除
	Created            time.Time `db:"created" json:"created"`
	CreatedBy          string    `db:"created_by" json:"created_by"`
	Updated            time.Time `db:"updated" json:"updated"`
	UpdatedBy          string    `db:"updated_by" json:"updated_by"`
}

type PaymentRecord struct {
	ID               int64     `db:"id" json:"id"`
	OrganizationID   int64     `db:"organization_id" json:"organization_id"`
	ProjectID        int64     `db:"project_id" json:"project_id"`
	PaymentRequestID int64     `db:"payment_request_id" json:"payment_request_id"`
	Date             time.Time `db:"date" json:"date"`
	Amount           float64   `db:"amount" json:"amount"`
	Remark           string    `db:"remark" json:"remark"`
	Picture          string    `db:"picture" json:"picture"`
	Status           int       `db:"status" json:"status"`
	Created          time.Time `db:"created" json:"created"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	Updated          time.Time `db:"updated" json:"updated"`
	UpdatedBy        string    `db:"updated_by" json:"updated_by"`
}

type PaymentRequestHistory struct {
	ID               int64     `db:"id" json:"id"`
	PaymentRequestID int64     `db:"payment_request_id" json:"payment_request_id"`
	Action           string    `db:"action" json:"action"`
	Content          string    `db:"content" json:"content"`
	Remark           string    `db:"remark" json:"remark"`
	Status           int       `db:"status" json:"status"`
	Created          time.Time `db:"created" json:"created"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	Updated          time.Time `db:"updated" json:"updated"`
	UpdatedBy        string    `db:"updated_by" json:"updated_by"`
}
