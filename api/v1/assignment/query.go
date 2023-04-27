package assignment

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type assignmentQuery struct {
	conn *sqlx.DB
}

func NewAssignmentQuery(connection *sqlx.DB) *assignmentQuery {
	return &assignmentQuery{
		conn: connection,
	}
}

func (r *assignmentQuery) GetAssignmentByID(id int64, organizationID int64) (*AssignmentResponse, error) {
	var assignment AssignmentResponse
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&assignment, `
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
		AND m.organization_id = ? 
		AND m.status > 0
		`, id, organizationID)
	} else {
		err = r.conn.Get(&assignment, `		
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
		m.created,
		m.created_by
		FROM assignments m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		LEFT JOIN projects p
		ON p.id = m.project_id
		LEFT JOIN users u
		ON u.id = m.assign_to,
		LEFT JOIN users u2
		ON u2.id = m.audit_to
		WHERE m.id = ? 
		AND m.status > 0
		`, id)
	}
	return &assignment, err
}

func (r *assignmentQuery) GetAssignmentCount(filter AssignmentFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Status; v != "all" {
		where, args = append(where, "status < ?"), append(args, 9)
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.AssignmentType; v != 0 {
		where, args = append(where, "assignment_type = ?"), append(args, v)
	}
	if v := filter.ReferenceID; v != 0 {
		where, args = append(where, "reference_id = ?"), append(args, v)
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM assignments 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *assignmentQuery) GetAssignmentList(filter AssignmentFilter) (*[]AssignmentResponse, error) {
	where, args := []string{"m.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "m.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Status; v != "all" {
		where, args = append(where, "m.status < ?"), append(args, 9)
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "m.organization_id = ?"), append(args, v)
	}
	if v := filter.AssignmentType; v != 0 {
		where, args = append(where, "m.assignment_type = ?"), append(args, v)
	}
	if v := filter.ReferenceID; v != 0 {
		where, args = append(where, "m.reference_id = ?"), append(args, v)
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "m.project_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var assignments []AssignmentResponse
	err := r.conn.Select(&assignments, `
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
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &assignments, nil
}

func (r *assignmentQuery) GetMyAssignmentCount(filter MyAssignmentFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Status; v != "all" {
		where = append(where, "status in (1,3)")
	}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "assign_to = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM assignments 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *assignmentQuery) GetMyAssignmentList(filter MyAssignmentFilter) (*[]AssignmentResponse, error) {
	where, args := []string{"m.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "m.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Status; v != "all" {
		where = append(where, "m.status in (1,3)")
	}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "m.assign_to = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var assignments []AssignmentResponse
	err := r.conn.Select(&assignments, `
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
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &assignments, nil
}

func (r *assignmentQuery) GetMyAuditCount(filter MyAuditFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Status; v != "all" {
		where = append(where, "status = 2")
	}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "audit_to = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM assignments 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *assignmentQuery) GetMyAuditList(filter MyAuditFilter) (*[]AssignmentResponse, error) {
	where, args := []string{"m.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "m.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Status; v != "all" {
		where = append(where, "m.status = 2")
	}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "m.audit_to = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var assignments []AssignmentResponse
	err := r.conn.Select(&assignments, `
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
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &assignments, nil
}
