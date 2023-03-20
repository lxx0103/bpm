package assignment

import "time"

type AssignmentFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	AssignmentType int    `form:"assignment_type" binding:"omitempty,min=1"`
	ReferenceID    int64  `form:"reference_id" binding:"omitempty,min=1"`
	ProjectID      int64  `form:"project_id" binding:"omitempty"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type AssignmentNew struct {
	OrganizationID int64  `json:"organization_id" binding:"omitempty,min=1"`
	AssignmentType int    `json:"assignment_type" binding:"required,oneof=1,2"`
	ReferenceID    int64  `json:"reference_id" binding:"omitempty,min=1"`
	ProjectID      int64  `json:"project_id" binding:"required"`
	AssignTo       int64  `json:"assign_to" binding:"required"`
	AuditTo        int64  `json:"audit_to" binding:"required"`
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Content        string `json:"content" binding:"required"`
	File           string `json:"file" binding:"omitempty"`
	User           string `json:"user" swaggerignore:"true"`
}

type AssignmentID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type AssignmentResponse struct {
	ID               int64     `db:"id" json:"id"`
	OrganizationID   int64     `db:"organization_id" json:"organization_id"`
	OrganizationName string    `db:"organization_name" json:"organization_name"`
	Date             time.Time `db:"date" json:"date"`
	Name             string    `db:"name" json:"name"`
	Content          string    `db:"content" json:"content"`
	File             string    `db:"file" json:"file"`
	Status           int       `db:"status" json:"status"`
}
