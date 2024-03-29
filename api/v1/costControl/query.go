package costControl

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type costControlQuery struct {
	conn *sqlx.DB
}

func NewCostControlQuery(connection *sqlx.DB) *costControlQuery {
	return &costControlQuery{
		conn: connection,
	}
}

func (q *costControlQuery) GetBudgetCount(filter ReqBudgetFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.BudgetType; v > 0 {
		where, args = append(where, "budget_type = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name LIKE ?"), append(args, "%"+v+"%")
	}
	var count int
	err := q.conn.Get(&count, `
		SELECT COUNT(*) 
		FROM budgets
		WHERE `+strings.Join(where, " AND "), args...)

	return count, err
}

func (q *costControlQuery) GetBudgetList(filter ReqBudgetFilter) (*[]RespBudget, error) {
	where, args := []string{"b.status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "b.project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "b.organization_id = ?"), append(args, v)
	}
	if v := filter.BudgetType; v > 0 {
		where, args = append(where, "b.budget_type = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "b.name LIKE ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var budgets []RespBudget
	err := q.conn.Select(&budgets, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	p.name AS project_name, 
	b.budget_type AS budget_type, 
	b.name AS name, 
	b.quantity AS quantity,
	b.unit_price AS unit_price,
	b.budget AS budget,
	b.used AS used,
	b.balance AS balance,
	b.remark AS remark,
	b.status AS status
	FROM budgets b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE `+strings.Join(where, " AND ")+`
	ORDER BY b.id DESC
	LIMIT ?, ?
	`, args...)
	return &budgets, err
}

func (q *costControlQuery) GetBudgetByID(id int64) (*RespBudget, error) {
	var budget RespBudget
	err := q.conn.Get(&budget, `
	SELECT
	b.id AS id,
	b.organization_id AS organization_id,	
	o.name AS organization_name,
	b.project_id AS project_id,
	p.name AS project_name, 
	b.budget_type AS budget_type, 
	b.name AS name, 
	b.quantity AS quantity,
	b.unit_price AS unit_price,
	b.budget AS budget,
	b.used AS used,
	b.balance AS balance,
	b.remark AS remark,
	b.status AS status
	FROM budgets b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE b.id = ? AND b.status = 1 
	`, id)
	return &budget, err
}

func (q *costControlQuery) GetBudgetPictureList(id int64) (*[]string, error) {
	var pictures []string
	err := q.conn.Select(&pictures, `
	SELECT link 
	FROM budget_pictures 
	WHERE budget_id = ? AND status = 1
	`, id)
	return &pictures, err
}

func (q *costControlQuery) GetPaymentRequestCount(filter ReqPaymentRequestFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentRequestType; v > 0 {
		where, args = append(where, "payment_request_type = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name LIKE ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v == "mine" {
		where, args = append(where, "user_id = ?"), append(args, filter.UserID)
	}
	if v := filter.Type; v == "passed" {
		where = append(where, "status in (2, 4, 5)")
	}
	if v := filter.Type; v == "audit" {
		where, args = append(where, "id in (SELECT p.id FROM payment_requests p LEFT JOIN payment_request_audits pa ON p.id = pa.payment_request_id AND p.audit_level = pa.audit_level WHERE p.status = 1 and pa.status > 0 AND p.organization_id = ? AND ( ( audit_type = 1 AND audit_to = ? ) OR ( audit_type = 2 and audit_to = ? ) ) )"), append(args, filter.OrganizationID, filter.PositionID, filter.UserID)
	}
	if v := filter.PaymentStatus; v == "none" {
		where = append(where, "paid = 0")
	}
	if v := filter.PaymentStatus; v == "partial" {
		where = append(where, "paid > 0 AND due > 0")
	}
	if v := filter.PaymentStatus; v == "paid" {
		where = append(where, "due = 0")
	}
	if v := filter.DeliveryStatus; v == "none" {
		where = append(where, "deliveried = 0")
	}
	if v := filter.DeliveryStatus; v == "partial" {
		where = append(where, "deliveried > 0 AND pending > 0")
	}
	if v := filter.DeliveryStatus; v == "deliveried" {
		where = append(where, "pending = 0")
	}
	var count int
	err := q.conn.Get(&count, `
		SELECT COUNT(*) 
		FROM payment_requests
		WHERE `+strings.Join(where, " AND "), args...)

	return count, err
}

func (q *costControlQuery) GetPaymentRequestList(filter ReqPaymentRequestFilter) (*[]RespPaymentRequest, error) {
	where, args := []string{"b.status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "b.project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "b.organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentRequestType; v > 0 {
		where, args = append(where, "b.payment_request_type = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "b.name LIKE ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v == "mine" {
		where, args = append(where, "b.user_id = ?"), append(args, filter.UserID)
	}
	if v := filter.Type; v == "passed" {
		where = append(where, "b.status in (2, 4, 5)")
	}
	if v := filter.Type; v == "audit" {
		where, args = append(where, "b.id in (SELECT p.id FROM payment_requests p LEFT JOIN payment_request_audits pa ON p.id = pa.payment_request_id AND p.audit_level = pa.audit_level WHERE p.status = 1 and pa.status > 0 AND p.organization_id = ? AND ( ( audit_type = 1 AND audit_to = ? ) OR ( audit_type = 2 and audit_to = ? ) ) )"), append(args, filter.OrganizationID, filter.PositionID, filter.UserID)
	}
	if v := filter.PaymentStatus; v == "none" {
		where = append(where, "b.paid = 0")
	}
	if v := filter.PaymentStatus; v == "partial" {
		where = append(where, "b.paid > 0 AND b.due > 0")
	}
	if v := filter.PaymentStatus; v == "paid" {
		where = append(where, "b.due = 0")
	}
	if v := filter.DeliveryStatus; v == "none" {
		where = append(where, "b.deliveried = 0")
	}
	if v := filter.DeliveryStatus; v == "partial" {
		where = append(where, "b.deliveried > 0 AND b.pending > 0")
	}
	if v := filter.DeliveryStatus; v == "deliveried" {
		where = append(where, "b.pending = 0")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var payment_requests []RespPaymentRequest
	err := q.conn.Select(&payment_requests, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	IFNULL(p.name, "") AS project_name, 
	b.payment_request_type AS payment_request_type, 
	b.budget_id AS budget_id,
	b.name AS name, 
	b.quantity AS quantity,
	b.unit_price AS unit_price,
	b.total AS total,
	b.paid AS paid,
	b.due AS due,
	b.remark AS remark,
	b.audit_level AS audit_level,
	b.status AS status,
	b.deliveried as deliveried,
	b.pending as pending,
	b.delivery_status as delivery_status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM payment_requests b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE `+strings.Join(where, " AND ")+`
	ORDER BY b.id DESC
	LIMIT ?, ?
	`, args...)
	return &payment_requests, err
}

func (q *costControlQuery) GetPaymentRequestByID(id int64) (*RespPaymentRequest, error) {
	var payment_request RespPaymentRequest
	err := q.conn.Get(&payment_request, `
	SELECT
	b.id AS id,
	b.organization_id AS organization_id,	
	o.name AS organization_name,
	b.project_id AS project_id,
	IFNULL(p.name, "") AS project_name, 
	b.payment_request_type AS payment_request_type, 
	b.budget_id AS budget_id,
	b.name AS name, 
	b.quantity AS quantity,
	b.unit_price AS unit_price,
	b.total AS total,
	b.paid AS paid,
	b.due AS due,
	b.remark AS remark,
	b.audit_level AS audit_level,
	b.status AS status,
	b.deliveried as deliveried,
	b.pending as pending,
	b.delivery_status as delivery_status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM payment_requests b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE b.id = ? AND b.status > 0
	`, id)
	return &payment_request, err
}

func (q *costControlQuery) GetPaymentRequestPictureList(id int64) (*[]string, error) {
	var pictures []string
	err := q.conn.Select(&pictures, `
	SELECT link 
	FROM payment_request_pictures 
	WHERE payment_request_id = ? AND status = 1
	`, id)
	return &pictures, err
}

func (q *costControlQuery) GetPaymentRequestTypeList(organizationID, paymentRequestType int64) (*[]RespPaymentRequestTypeAudit, error) {
	var res []RespPaymentRequestTypeAudit
	err := q.conn.Select(&res, `
		SELECT pr.audit_level AS audit_level,
		pr.audit_type AS audit_type,
		pr.audit_to AS audit_to,
		CASE WHEN pr.audit_type = 1 THEN IFNULL(p.name,"") ELSE IFNULL(u.name,"") END AS audit_to_name		
		FROM payment_request_type_audits pr
		LEFT JOIN positions p
		ON pr.audit_to = p.id
		LEFT JOIN users u
		ON pr.audit_to = u.id
		WHERE pr.organization_id = ?
		AND pr.payment_request_type = ?
		AND pr.status = 1
		ORDER BY pr.audit_level ASC
	`, organizationID, paymentRequestType)
	return &res, err
}

func (q *costControlQuery) GetPaymentRequestHistoryList(paymentRequestID int64) (*[]RespPaymentRequestHistory, error) {
	var res []RespPaymentRequestHistory
	err := q.conn.Select(&res, `
		SELECT id, payment_request_id, action, content, remark, created_by, created 
		FROM payment_request_historys
		WHERE payment_request_id = ?
		AND status > 0
		ORDER BY id desc
	`, paymentRequestID)
	return &res, err
}

func (q *costControlQuery) GetPaymentRequestHistoryPictureList(id int64) (*[]string, error) {
	var pictures []string
	err := q.conn.Select(&pictures, `
	SELECT link 
	FROM payment_request_history_pictures 
	WHERE payment_request_history_id = ? AND status = 1
	`, id)
	return &pictures, err
}

func (q *costControlQuery) GetPaymentRequestAuditList(paymentRequestID int64) (*[]RespPaymentRequestAudit, error) {
	var res []RespPaymentRequestAudit
	err := q.conn.Select(&res, `
		SELECT pr.audit_level AS audit_level,
		pr.audit_type AS audit_type,
		pr.audit_to AS audit_to,
		CASE WHEN pr.audit_type = 1 THEN IFNULL(p.name,"") ELSE IFNULL(u.name,"") END AS audit_to_name		
		FROM payment_request_audits pr
		LEFT JOIN positions p
		ON pr.audit_to = p.id
		LEFT JOIN users u
		ON pr.audit_to = u.id
		WHERE pr.payment_request_id = ?
		AND pr.status = 1
		ORDER BY pr.audit_level ASC
	`, paymentRequestID)
	return &res, err
}

func (q *costControlQuery) GetPaymentCount(filter ReqPaymentFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentRequestID; v > 0 {
		where, args = append(where, "payment_request_id = ?"), append(args, v)
	}
	var count int
	err := q.conn.Get(&count, `
		SELECT COUNT(*) 
		FROM payments
		WHERE `+strings.Join(where, " AND "), args...)

	return count, err
}

func (q *costControlQuery) GetPaymentList(filter ReqPaymentFilter) (*[]RespPayment, error) {
	where, args := []string{"b.status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "b.project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "b.organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentRequestID; v > 0 {
		where, args = append(where, "b.payment_request_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var payments []RespPayment
	err := q.conn.Select(&payments, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	IFNULL(p.name, "") AS project_name, 
	b.payment_request_id AS payment_request_id, 
	b.payment_date AS payment_date,
	b.amount AS amount, 
	b.payment_method AS payment_method,
	b.remark AS remark,
	b.user_id AS user_id,
	b.status AS status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM payments b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE `+strings.Join(where, " AND ")+`
	ORDER BY b.id DESC
	LIMIT ?, ?
	`, args...)
	return &payments, err
}

func (q *costControlQuery) GetPaymentByID(id int64) (*RespPayment, error) {
	var payment RespPayment
	err := q.conn.Get(&payment, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	IFNULL(p.name, "") AS project_name, 
	b.payment_request_id AS payment_request_id, 
	b.payment_date AS payment_date,
	b.amount AS amount, 
	b.payment_method AS payment_method,
	b.remark AS remark,
	b.user_id AS user_id,
	b.status AS status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM payments b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE b.id = ? AND b.status > 0
	`, id)
	return &payment, err
}

func (q *costControlQuery) GetPaymentPictureList(id int64) (*[]string, error) {
	var pictures []string
	err := q.conn.Select(&pictures, `
	SELECT link 
	FROM payment_pictures 
	WHERE payment_id = ? AND status = 1
	`, id)
	return &pictures, err
}

func (q *costControlQuery) GetIncomeCount(filter ReqIncomeFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	var count int
	err := q.conn.Get(&count, `
		SELECT COUNT(*) 
		FROM incomes
		WHERE `+strings.Join(where, " AND "), args...)

	return count, err
}

func (q *costControlQuery) GetIncomeList(filter ReqIncomeFilter) (*[]RespIncome, error) {
	where, args := []string{"b.status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "b.project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "b.organization_id = ?"), append(args, v)
	}

	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var incomes []RespIncome
	err := q.conn.Select(&incomes, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	IFNULL(p.name, "") AS project_name, 
	b.title AS title, 
	b.date AS date,
	b.amount AS amount, 
	b.payment_method AS payment_method,
	b.remark AS remark,
	b.user_id AS user_id,
	b.status AS status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM incomes b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE `+strings.Join(where, " AND ")+`
	ORDER BY b.id DESC
	LIMIT ?, ?
	`, args...)
	return &incomes, err
}

func (q *costControlQuery) GetIncomeByID(id int64) (*RespIncome, error) {
	var income RespIncome
	err := q.conn.Get(&income, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	IFNULL(p.name, "") AS project_name, 
	b.title AS title, 
	b.date AS date,
	b.amount AS amount, 
	b.payment_method AS payment_method,
	b.remark AS remark,
	b.user_id AS user_id,
	b.status AS status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM incomes b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE b.id = ? AND b.status > 0
	`, id)
	return &income, err
}

func (q *costControlQuery) GetIncomePictureList(id int64) (*[]string, error) {
	var pictures []string
	err := q.conn.Select(&pictures, `
	SELECT link 
	FROM income_pictures 
	WHERE income_id = ? AND status = 1
	`, id)
	return &pictures, err
}

func (q *costControlQuery) GetDeliveryCount(filter ReqDeliveryFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentRequestID; v > 0 {
		where, args = append(where, "payment_request_id = ?"), append(args, v)
	}
	var count int
	err := q.conn.Get(&count, `
		SELECT COUNT(*) 
		FROM deliverys
		WHERE `+strings.Join(where, " AND "), args...)

	return count, err
}

func (q *costControlQuery) GetDeliveryList(filter ReqDeliveryFilter) (*[]RespDelivery, error) {
	where, args := []string{"b.status > 0"}, []interface{}{}
	if v := filter.ProjectID; v > 0 {
		where, args = append(where, "b.project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v > 0 {
		where, args = append(where, "b.organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentRequestID; v > 0 {
		where, args = append(where, "b.payment_request_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var payments []RespDelivery
	err := q.conn.Select(&payments, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	IFNULL(p.name, "") AS project_name, 
	b.payment_request_id AS payment_request_id, 
	b.date AS date,
	b.quantity AS quantity, 
	b.remark AS remark,
	b.user_id AS user_id,
	b.status AS status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM deliverys b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE `+strings.Join(where, " AND ")+`
	ORDER BY b.id DESC
	LIMIT ?, ?
	`, args...)
	return &payments, err
}

func (q *costControlQuery) GetDeliveryByID(id int64) (*RespDelivery, error) {
	var payment RespDelivery
	err := q.conn.Get(&payment, `
	SELECT b.id AS id, 
	b.organization_id AS organization_id, 
	o.name AS organization_name, 
	b.project_id AS project_id, 
	IFNULL(p.name, "") AS project_name, 
	b.payment_request_id AS payment_request_id, 
	b.date AS date,
	b.quantity AS quantity, 
	b.remark AS remark,
	b.user_id AS user_id,
	b.status AS status,
	b.user_id AS user_id,
	b.created AS created,
	b.created_by AS created_by
	FROM deliverys b
	LEFT JOIN projects p ON b.project_id = p.id
	LEFT JOIN organizations o ON b.organization_id = o.id
	WHERE b.id = ? AND b.status > 0
	`, id)
	return &payment, err
}

func (q *costControlQuery) GetDeliveryPictureList(id int64) (*[]string, error) {
	var pictures []string
	err := q.conn.Select(&pictures, `
	SELECT link 
	FROM delivery_pictures 
	WHERE delivery_id = ? AND status = 1
	`, id)
	return &pictures, err
}

func (q *costControlQuery) GetBudgetSumByProjectID(projectID int64) (float64, error) {
	var sum float64
	err := q.conn.Get(&sum, `
	SELECT IFNULL(sum(budget), 0)
	FROM budgets 
	WHERE project_id = ? AND status = 1
	`, projectID)
	return sum, err
}

func (q *costControlQuery) GetIncomeSumByProjectID(projectID int64) (float64, error) {
	var sum float64
	err := q.conn.Get(&sum, `
	SELECT IFNULL(sum(amount), 0)
	FROM incomes 
	WHERE project_id = ? AND status = 1
	`, projectID)
	return sum, err
}

func (q *costControlQuery) GetPaymentSumByProjectID(projectID int64) (float64, error) {
	var sum float64
	err := q.conn.Get(&sum, `
	SELECT IFNULL(sum(amount), 0)
	FROM payments 
	WHERE project_id = ? AND status = 1
	`, projectID)
	return sum, err
}
