package assignment

import "time"

type Assignment struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	AssignmentType  int       `db:"assignment_type" json:"assignment_type"`
	ReferenceID     int64     `db:"reference_id" json:"reference_id"`
	AssignTo        int64     `db:"assign_to" json:"assign_to"`
	AuditTo         int64     `db:"audit_to" json:"audit_to"`
	CompleteContent string    `db:"complete_content" json:"complete_content"`
	AuditContent    string    `db:"audit_content" json:"audit_content"`
	Name            string    `db:"name" json:"name"`
	Content         string    `db:"content" json:"content"`
	File            string    `db:"file" json:"file"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}
