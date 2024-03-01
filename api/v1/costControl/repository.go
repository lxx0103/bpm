package costControl

import (
	"database/sql"
	"fmt"
	"time"
)

type costControlRepository struct {
	tx *sql.Tx
}

func NewCostControlRepository(transaction *sql.Tx) *costControlRepository {
	return &costControlRepository{
		tx: transaction,
	}
}

func (r *costControlRepository) CreateBudget(info ReqBudgetNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO budgets 
		(
			organization_id,
			project_id,
			budget_type,
			name,
			quantity,
			unit_price,
			budget,
			used,
			balance,
			remark,
			status,
			created,
			created_by,
			updated,
			updated_by
		) 
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`, info.OrganizationID, info.ProjectID, info.BudgetType, info.Name, info.Quantity, info.UnitPrice, info.Budget, 0, info.Budget, info.Remark, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *costControlRepository) CreateBudgetPicture(info ReqBudgetPictureNew) error {
	_, err := r.tx.Exec(`
	INSERT INTO budget_pictures 
	(
		budget_id,
		link,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?
	)`, info.BudgetID, info.Picture, 1, time.Now(), info.User, time.Now(), info.User)

	return err
}

func (r *costControlRepository) DeleteBudgetPicture(budgetID int64) error {
	_, err := r.tx.Exec(`
	UPDATE budget_pictures 
	SET status = -1 
	WHERE budget_id = ?
	`, budgetID)
	return err
}

func (r *costControlRepository) UpdateBudget(info ReqBudgetUpdate, id int64) error {
	_, err := r.tx.Exec(`
		UPDATE budgets SET 
			name = ?,
			quantity = ?,
			unit_price = ?,
			budget = ?,
			used = ?,
			balance = ?,
			remark = ?,
			updated = ?,
			updated_by = ?
		WHERE id = ?
	`, info.Name, info.Quantity, info.UnitPrice, info.Budget, info.Used, info.Balance, info.Remark, time.Now(), info.User, id)
	return err
}

func (r *costControlRepository) GetBudgetByID(id int64) (RespBudget, error) {
	var budget RespBudget
	row := r.tx.QueryRow(`
		SELECT id,
		organization_id,
		project_id,
		budget_type,
		name,
		quantity,
		unit_price,
		budget,
		used,
		balance,
		remark,
		status
		FROM budgets
		WHERE id = ?
		AND status = 1
	`, id)
	err := row.Scan(&budget.ID, &budget.OrganizationID, &budget.ProjectID, &budget.BudgetType, &budget.Name, &budget.Quantity, &budget.UnitPrice, &budget.Budget, &budget.Used, &budget.Balance, &budget.Remark, &budget.Status)
	return budget, err
}

func (r *costControlRepository) DeleteBudget(id int64, user string) error {
	_, err := r.tx.Exec(`
		UPDATE budgets SET 
			status = -1,
			updated = ?,
			updated_by = ?
		WHERE id = ?
	`, time.Now(), user, id)
	return err
}

func (r *costControlRepository) CreatePaymentRequest(info ReqPaymentRequestNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO payment_requests 
		(
			organization_id,
			project_id,
			budget_id,
			payment_request_type,
			name,
			quantity,
			unit_price,
			total,
			paid,
			due,
			status,
			remark,
			audit_level,
			user_id,
			created,
			created_by,
			updated,
			updated_by
		) 
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`, info.OrganizationID, info.ProjectID, info.BudgetID, info.PaymentRequestType, info.Name, info.Quantity, info.UnitPrice, info.Total, 0, info.Total, 1, info.Remark, 1, info.UserID, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (r *costControlRepository) CreatePaymentRequestPicture(info ReqPaymentRequestPictureNew) error {
	_, err := r.tx.Exec(`
	INSERT INTO payment_request_pictures 
	(
		payment_request_id,
		link,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?
	)`, info.PaymentRequestID, info.Picture, 1, time.Now(), info.User, time.Now(), info.User)

	return err

}

func (r *costControlRepository) DeletePaymentRequestPicture(paymentRequestID int64) error {
	_, err := r.tx.Exec(`
	UPDATE payment_request_pictures 
	SET status = -1 
	WHERE payment_request_id = ?
	`, paymentRequestID)
	return err
}

func (r *costControlRepository) CreatePaymentRequestHistory(info ReqPaymentRequestHistoryNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO payment_request_historys 
		(
			payment_request_id,
			action,
			content,
			remark,
			status,
			created,
			created_by,
			updated,
			updated_by
		) 
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		)`, info.PaymentRequestID, info.Action, info.Content, info.Remark, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *costControlRepository) GetPaymentRequestByID(id int64) (RespPaymentRequest, error) {
	var paymentRequest RespPaymentRequest
	row := r.tx.QueryRow(`
	    SELECT organization_id,
		project_id,
		budget_id,
		payment_request_type,
		name,
		quantity,
		unit_price,
		total,
		paid,
		due,
		remark,
		audit_level,
		user_id,
		status
		FROM payment_requests
		WHERE id = ?
		AND status > 0
	`, id)
	err := row.Scan(&paymentRequest.OrganizationID, &paymentRequest.ProjectID, &paymentRequest.BudgetID, &paymentRequest.PaymentRequestType, &paymentRequest.Name, &paymentRequest.Quantity, &paymentRequest.UnitPrice, &paymentRequest.Total, &paymentRequest.Paid, &paymentRequest.Due, &paymentRequest.Remark, &paymentRequest.AuditLevel, &paymentRequest.UserID, &paymentRequest.Status)
	return paymentRequest, err
}

func (r *costControlRepository) UpdatePaymentRequest(info ReqPaymentRequestUpdate, id int64) error {
	_, err := r.tx.Exec(`
		UPDATE payment_requests SET 
			project_id = ?,
			payment_request_type = ?,
			budget_id = ?,
			name = ?,
			quantity = ?,
			unit_price = ?,
			total = ?,
			paid = ?,
			due = ?,
			remark = ?,
			audit_level = ?,
			status = ?,
			updated = ?,
			updated_by = ?
		WHERE id = ?
	`, info.ProjectID, info.PaymentRequestType, info.BudgetID, info.Name, info.Quantity, info.UnitPrice, info.Total, 0, info.Total, info.Remark, 1, info.Status, time.Now(), info.User, id)
	return err
}

func (r *costControlRepository) DeletePaymentRequest(id int64, user string) error {
	_, err := r.tx.Exec(`
		UPDATE payment_requests SET 
			status = -1,
			updated = ?,
			updated_by = ?
		WHERE id = ?
	`, time.Now(), user, id)
	return err
}

func (r *costControlRepository) DeletePaymentRequestTypeAudit(paymentRequestType int, organizationID int64, byUser string) error {
	_, err := r.tx.Exec(`
	    UPDATE payment_request_type_audits SET
			status = -1,
			updated = ?,
			updated_by = ?
		WHERE organization_id = ?
		AND payment_request_type = ?	
	`, time.Now(), byUser, organizationID, paymentRequestType)
	return err
}

func (r *costControlRepository) CreatePaymentRequestTypeAudit(info ReqPaymentRequestTypeAudit) error {
	_, err := r.tx.Exec(`
		INSERT INTO payment_request_type_audits 
		(
			organization_id,
			payment_request_type,
			audit_level,
			audit_type,
			audit_to,
			status,
			created,
			created_by,
			updated,
			updated_by
		) VALUES
		(
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`, info.OrganizationID, info.PaymentRequestType, info.AuditLevel, info.AuditType, info.AuditTo, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *costControlRepository) GetPaymentRequestTypeAudit(organizationID int64, paymentRequestType int) (*[]RespPaymentRequestTypeAudit, error) {
	var res []RespPaymentRequestTypeAudit
	rows, err := r.tx.Query(`
		SELECT audit_level, audit_type, audit_to 
		FROM payment_request_type_audits 
		WHERE organization_id = ?
		AND payment_request_type = ?
		AND status = 1
	`, organizationID, paymentRequestType)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes RespPaymentRequestTypeAudit
		err = rows.Scan(&rowRes.AuditLevel, &rowRes.AuditType, &rowRes.AuditTo)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *costControlRepository) CreatePaymentRequestAudit(info ReqPaymentRequestAuditNew) error {
	_, err := r.tx.Exec(`
	INSERT INTO payment_request_audits 
	(
		payment_request_id,
		audit_level,
		audit_type,
		audit_to,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?
	)`, info.PaymentRequestID, info.AuditLevel, info.AuditType, info.AuditTo, 1, time.Now(), info.User, time.Now(), info.User)

	return err

}

func (r *costControlRepository) DeletePaymentRequestAudits(paymentRequestID int64) error {
	_, err := r.tx.Exec(`
	UPDATE payment_request_audits 
	SET status = -1 
	WHERE payment_request_id = ?
	`, paymentRequestID)
	return err
}

func (r *costControlRepository) CheckAudit(paymentRequestID int64, userID int64, positionID int64, auditLevel int) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM payment_request_audits WHERE payment_request_id = ? AND ( ( audit_type = 1 AND audit_to = ? ) OR ( audit_type = 2 and audit_to = ? ) ) AND status > 0 AND audit_level = ? LIMIT 1`, paymentRequestID, positionID, userID, auditLevel)
	err := row.Scan(&res)
	return res, err
}

func (r *costControlRepository) GetNextLevel(paymentRequestID int64, auditLevel int) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT audit_level FROM payment_request_audits WHERE payment_request_id = ? AND status > 0 AND audit_level > ? ORDER by audit_level ASC LIMIT 1`, paymentRequestID, auditLevel)
	err := row.Scan(&res)
	return res, err
}

func (r *costControlRepository) AuditPaymentRequest(paymentRequestID int64, auditLevel, status int, byUser string) error {
	_, err := r.tx.Exec(`
	UPDATE payment_requests 
	SET audit_level = ?,
	status = ?,
	updated = ?,
	updated_by = ?
	WHERE id = ?
	`, auditLevel, status, time.Now(), byUser, paymentRequestID)
	return err
}

func (r *costControlRepository) CreatePaymentRequestHistoryPicture(info ReqPaymentRequestHistoryPictureNew) error {
	_, err := r.tx.Exec(`
	INSERT INTO payment_request_history_pictures 
	(
		payment_request_history_id,
		link,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?
	)`, info.PaymentRequestHistoryID, info.Picture, 1, time.Now(), info.User, time.Now(), info.User)

	return err

}

func (r *costControlRepository) DeletePaymentRequestAudit(paymentRequestID int64, organizationID int64, byUser string) error {
	_, err := r.tx.Exec(`
	    UPDATE payment_request_audits SET
			status = -1,
			updated = ?,
			updated_by = ?
		WHERE payment_request_id = ?	
	`, time.Now(), byUser, paymentRequestID)
	return err
}

func (r *costControlRepository) CreatePayment(paymentRequestID int64, info ReqPaymentNew) (int64, error) {
	result, err := r.tx.Exec(`
	INSERT INTO payments 
	(
		organization_id,
		project_id,
		payment_request_id,
		amount,
		payment_method,
		payment_date,
		remark,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	)`, info.OrganizationID, info.ProjectID, paymentRequestID, info.Amount, info.PaymentMethod, info.PaymentDate, info.Remark, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return result.LastInsertId()
}

func (r *costControlRepository) CreatePaymentPicture(info ReqPaymentPictureNew) error {
	_, err := r.tx.Exec(`
	INSERT INTO payment_pictures 
	(
		payment_id,
		link,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?
	)`, info.PaymentID, info.Picture, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *costControlRepository) UpdatePaymentRequestPaid(paymentRequestID int64, info ReqPaymentRequestPaid) error {
	_, err := r.tx.Exec(`
		UPDATE payment_requests SET 
			paid = ?,
			due = ?,
			status = ?,
			updated = ?,
			updated_by = ?
		WHERE id = ?
	`, info.Paid, info.Due, info.Status, time.Now(), info.User, paymentRequestID)
	return err
}

func (r *costControlRepository) UpdateBudgitUsed(budgetID int64, info ReqBudgetPaid) error {
	_, err := r.tx.Exec(`
		UPDATE budgets SET 
			used = ?,
			balance = ?,
			updated = ?,
			updated_by = ?
		WHERE id = ?
	`, info.Used, info.Balance, time.Now(), info.User, budgetID)
	return err
}

func (r *costControlRepository) GetPaymentByID(id int64) (RespPayment, error) {
	var payment RespPayment
	row := r.tx.QueryRow(`
	    SELECT organization_id,
		project_id,
		payment_request_id,
		payment_date,
		amount,
		payment_method,
		remark,
		user_id,
		status
		FROM payments
		WHERE id = ?
		AND status > 0
	`, id)
	err := row.Scan(&payment.OrganizationID, &payment.ProjectID, &payment.PaymentRequestID, &payment.PaymentDate, &payment.Amount, &payment.PaymentMethod, &payment.Remark, &payment.UserID, &payment.Status)
	return payment, err
}

func (r *costControlRepository) DeletePaymentPicture(paymentID int64) error {
	_, err := r.tx.Exec(`
	UPDATE payment_pictures 
	SET status = -1 
	WHERE payment_id = ?
	`, paymentID)
	return err
}

func (r *costControlRepository) DeletePayment(id int64, user string) error {
	_, err := r.tx.Exec(`
	UPDATE payments 
	SET status = -1 
	updated = ?,
	updated_by = ?
	WHERE id = ?
	`, id, time.Now(), user)
	return err
}

func (r *costControlRepository) UpdatePayment(paymentID int64, info ReqPaymentUpdate) error {
	_, err := r.tx.Exec(`
		UPDATE payments SET 
			amount = ?,
			payment_method = ?,
			payment_date = ?,
			remark = ?,
			updated = ?,
			updated_by = ?
		WHERE id = ?
	`, info.Amount, info.PaymentMethod, info.PaymentDate, info.Remark, time.Now(), info.User, paymentID)
	return err
}

func (r *costControlRepository) CreateIncome(info ReqIncomeNew) (int64, error) {
	result, err := r.tx.Exec(`
	INSERT INTO incomes 
	(
		organization_id,
		project_id,
		title,
		amount,
		payment_method,
		date,
		remark,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	)`, info.OrganizationID, info.ProjectID, info.Title, info.Amount, info.PaymentMethod, info.Date, info.Remark, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return result.LastInsertId()
}

func (r *costControlRepository) CreateIncomePicture(info ReqIncomePictureNew) error {
	_, err := r.tx.Exec(`
	INSERT INTO income_pictures 
	(
		income_id,
		link,
		status,
		created,
		created_by,
		updated,
		updated_by
	) 
	VALUES (
		?, ?, ?, ?, ?, ?, ?
	)`, info.IncomeID, info.Picture, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}
