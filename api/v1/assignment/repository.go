package assignment

import (
	"database/sql"
	"time"
)

type assignmentRepository struct {
	tx *sql.Tx
}

func NewAssignmentRepository(transaction *sql.Tx) *assignmentRepository {
	return &assignmentRepository{
		tx: transaction,
	}
}

func (r *assignmentRepository) CreateAssignment(info AssignmentNew) error {
	_, err := r.tx.Exec(`
		INSERT INTO assignments
		(
			organization_id,
			project_id,
			assignment_type,
			reference_id,
			assign_to,
			audit_to,
			name,
			content,
			file,
			status,
			user_id,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ProjectID, info.AssignmentType, info.ReferenceID, info.AssignTo, info.AuditTo, info.Name, info.Content, info.File, 1, info.UserID, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *assignmentRepository) UpdateAssignment(id int64, info AssignmentUpdate) error {
	_, err := r.tx.Exec(`
		Update assignments SET 
		project_id = ?,
		assign_to = ?,
		audit_to = ?,
		name = ?,
		content = ?,
		file = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.ProjectID, info.AssignTo, info.AuditTo, info.Name, info.Content, info.File, time.Now(), info.User, id)
	return err
}

func (r *assignmentRepository) GetAssignmentByID(id int64) (*AssignmentResponse, error) {
	var res AssignmentResponse
	row := r.tx.QueryRow(`
		SELECT 
		m.id,
		m.organization_id, 
		o.name as organization_name,
		m.assignment_type,
		m.reference_id,
		m.project_id,
		IFNULL(p.name, "") as project_name,
		m.assign_to,
		IFNULL(u.name, "") as assign_name,
		m.audit_to,
		IFNULL(u2.name, "") as audit_name,
		m.complete_content,
		m.complete_time,
		m.audit_content,
		m.audit_time,
		m.name,
		m.content, 
		m.file,
		m.status,
		m.user_id,
		m.created,
		m.created_by
		FROM assignments m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		LEFT JOIN projects p
		ON p.id = m.project_id
		LEFT JOIN users u
		ON u.id = m.assign_to
		LEFT JOIN users u2
		ON u2.id = m.audit_to
		WHERE m.id = ? 
		AND m.status > 0
	`, id)

	err := row.Scan(&res.ID, &res.OrganizationID, &res.OrganizationName, &res.AssignmentType, &res.ReferenceID, &res.ProjectID, &res.ProjectName, &res.AssignTo, &res.AssignName, &res.AuditTo, &res.AuditName, &res.CompleteContent, &res.CompleteTime, &res.AuditContent, &res.AuditTime, &res.Name, &res.Content, &res.File, &res.Status, &res.UserID, &res.Created, &res.CreatedBy)
	return &res, err
}

func (r *assignmentRepository) DeleteAssignment(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update assignments SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *assignmentRepository) CompleteAssignment(id int64, info AssignmentComplete) error {
	_, err := r.tx.Exec(`
		Update assignments SET 
		complete_content = ?,
		complete_time = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Content, time.Now(), 2, time.Now(), info.User, id)
	return err
}

func (r *assignmentRepository) AuditAssignment(id int64, info AssignmentAudit) error {
	status := 9
	if info.Result == 2 {
		status = 3
	}
	_, err := r.tx.Exec(`
		Update assignments SET 
		audit_content = ?,
		audit_time = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Content, time.Now(), status, time.Now(), info.User, id)
	return err
}
