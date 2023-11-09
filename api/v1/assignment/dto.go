package assignment

import "time"

type AssignmentFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	Status         string `form:"status" binding:"required,oneof=all active"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	AssignmentType int    `form:"assignment_type" binding:"omitempty,min=1"`
	ReferenceID    int64  `form:"reference_id" binding:"omitempty,min=1"`
	ProjectID      int64  `form:"project_id" binding:"omitempty"`
	EventID        int64  `form:"event_id" binding:"omitempty"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type AssignmentNew struct {
	OrganizationID int64    `json:"organization_id" binding:"omitempty,min=1"`
	AssignmentType int      `json:"assignment_type" binding:"required,oneof=1 2"`
	ReferenceID    int64    `json:"reference_id" binding:"omitempty,min=1"`
	ProjectID      int64    `json:"project_id" binding:"required"`
	EventID        int64    `json:"event_id" binding:"omitempty"`
	AssignTo       int64    `json:"assign_to" binding:"required"`
	AuditTo        int64    `json:"audit_to" binding:"required"`
	Name           string   `json:"name" binding:"required,min=1,max=64"`
	Content        string   `json:"content" binding:"required"`
	File           []string `json:"file" binding:"omitempty"`
	User           string   `json:"user" swaggerignore:"true"`
	UserID         int64    `json:"user_id" swaggerignore:"true"`
}

type AssignmentUpdate struct {
	ProjectID int64    `json:"project_id" binding:"required"`
	EventID   int64    `json:"event_id" binding:"omitempty"`
	AssignTo  int64    `json:"assign_to" binding:"required"`
	AuditTo   int64    `json:"audit_to" binding:"required"`
	Name      string   `json:"name" binding:"required,min=1,max=64"`
	Content   string   `json:"content" binding:"required"`
	File      []string `json:"file" binding:"omitempty"`
	User      string   `json:"user" swaggerignore:"true"`
	UserID    int64    `json:"user_id" swaggerignore:"true"`
}
type AssignmentID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type AssignmentResponse struct {
	ID               int64     `db:"id" json:"id"`
	OrganizationID   int64     `db:"organization_id" json:"organization_id"`
	OrganizationName string    `db:"organization_name" json:"organization_name"`
	AssignmentType   int       `db:"assignment_type" json:"assignment_type"`
	ReferenceID      int64     `db:"reference_id" json:"reference_id"`
	ProjectID        int64     `db:"project_id" json:"project_id"`
	ProjectName      string    `db:"project_name" json:"project_name"`
	EventID          int64     `db:"event_id" json:"event_id"`
	EventName        string    `db:"event_name" json:"event_name"`
	AssignTo         int64     `db:"assign_to" json:"assign_to"`
	AssignName       string    `db:"assign_name" json:"assign_name"`
	AuditTo          int64     `db:"audit_to" json:"audit_to"`
	AuditName        string    `db:"audit_name" json:"audit_name"`
	CompleteContent  string    `db:"complete_content" json:"complete_content"`
	CompleteTime     string    `db:"complete_time" json:"complete_time"`
	AuditContent     string    `db:"audit_content" json:"audit_content"`
	AuditTime        string    `db:"audit_time" json:"audit_time"`
	Name             string    `db:"name" json:"name"`
	Content          string    `db:"content" json:"content"`
	File             []string  `db:"file" json:"file"`
	CompleteFile     []string  `db:"complete_file" json:"complete_file"`
	Status           int       `db:"status" json:"status"`
	UserID           int64     `db:"user_id" json:"user_id"`
	Created          time.Time `db:"created" json:"created"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
}

type AssignmentComplete struct {
	Content string   `json:"content" binding:"required"`
	File    []string `json:"file" binding:"omitempty"`
	User    string   `json:"user" swaggerignore:"true"`
	UserID  int64    `json:"user_id" swaggerignore:"true"`
}

type AssignmentAudit struct {
	Result  int    `json:"result" binding:"required,oneof=1 2"`
	Content string `json:"content" binding:"omitempty"`
	User    string `json:"user" swaggerignore:"true"`
	UserID  int64  `json:"user_id" swaggerignore:"true"`
}

type MyAssignmentFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	Status   string `form:"status" binding:"required,oneof=all active"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type MyAuditFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	Status   string `form:"status" binding:"required,oneof=all active"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type AssignmentFile struct {
	ID           int64     `db:"id" json:"id"`
	AssignmentID int64     `db:"assignment_id" json:"assignment_id"`
	Link         string    `db:"link" json:"link"`
	Status       int       `db:"status" json:"status"`
	Created      time.Time `db:"created" json:"created"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	Updated      time.Time `db:"updated" json:"updated"`
	UpdatedBy    string    `db:"updated_by" json:"updated_by"`
}

type AssignmentCompleteFile struct {
	ID           int64     `db:"id" json:"id"`
	AssignmentID int64     `db:"assignment_id" json:"assignment_id"`
	Link         string    `db:"link" json:"link"`
	Status       int       `db:"status" json:"status"`
	Created      time.Time `db:"created" json:"created"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	Updated      time.Time `db:"updated" json:"updated"`
	UpdatedBy    string    `db:"updated_by" json:"updated_by"`
}
