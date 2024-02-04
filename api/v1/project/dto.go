package project

import "time"

type ProjectFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	Type           int    `form:"type" binding:"omitempty,oneof=1 2"`
	Priority       int    `form:"priority" binding:"omitempty,oneof=1 2 3"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ProjectNew struct {
	Name            string  `json:"name" binding:"required,min=1,max=64"`
	TemplateID      int64   `json:"template_id" binding:"required,min=1"`
	ClientID        int64   `json:"client_id" binding:"omitempty,min=1"`
	Type            int     `json:"type" swaggerignore:"true"`
	Location        string  `json:"location" binding:"omitempty,min=1,max=64"`
	Longitude       float64 `json:"longitude" binding:"omitempty"`
	Latitude        float64 `json:"latitude" binding:"omitempty"`
	CheckinDistance int     `json:"checkin_distance" binding:"omitempty"`
	Priority        int     `json:"priority" binding:"required,oneof=1 2 3"`
	TeamID          int64   `json:"team_id" binding:"omitempty"`
	Area            string  `json:"area" binding:"omitempty,min=1,max=64"`
	RecordAlertDay  int     `json:"record_alert_day" binding:"omitempty,min=1"`
	User            string  `json:"user" swaggerignore:"true"`
	UserID          int64   `json:"user_id" swaggerignore:"true"`
}

type ProjectID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ProjectUpdate struct {
	Name            string  `json:"name" binding:"omitempty,min=1,max=64"`
	ClientID        int64   `json:"client_id" binding:"omitempty,min=1"`
	Location        string  `json:"location" binding:"omitempty,min=1,max=64"`
	Longitude       float64 `json:"longitude" binding:"omitempty"`
	Latitude        float64 `json:"latitude" binding:"omitempty"`
	CheckinDistance int     `json:"checkin_distance" binding:"omitempty"`
	Priority        int     `json:"priority" binding:"omitempty,oneof=1 2 3"`
	TeamID          int64   `json:"team_id" binding:"omitempty"`
	Area            string  `json:"area" binding:"omitempty,min=1,max=64"`
	RecordAlertDay  int     `json:"record_alert_day" binding:"omitempty,min=1"`
	User            string  `json:"user" swaggerignore:"true"`
	UserID          int64   `json:"user_id" swaggerignore:"true"`
}
type MyProjectFilter struct {
	Status   string `form:"status" binding:"required,oneof=all active"`
	Type     int    `form:"type" binding:"omitempty,oneof=1 2"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}
type AssignedProjectFilter struct {
	Name     string `form:"name" binding:"omitempty"`
	Status   int    `form:"status" binding:"omitempty,oneof=1 2"`
	Type     int    `form:"type" binding:"omitempty,oneof=1 2"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ProjectResponse struct {
	ID               int64                 `db:"id" json:"id"`
	OrganizationID   int64                 `db:"organization_id" json:"organization_id"`
	OrganizationName string                `db:"organization_name" json:"organization_name"`
	TemplateID       int64                 `db:"template_id" json:"template_id"`
	TemplateName     string                `db:"template_name" json:"template_name"`
	ClientID         int64                 `db:"client_id" json:"client_id"`
	ClientName       string                `db:"client_name" json:"client_name"`
	Name             string                `db:"name" json:"name"`
	Type             int                   `db:"type" json:"type"`
	Location         string                `db:"location" json:"location"`
	Longitude        float64               `db:"longitude" json:"longitude"`
	Latitude         float64               `db:"latitude" json:"latitude"`
	CheckinDistance  int                   `db:"checkin_distance" json:"checkin_distance"`
	Priority         int                   `db:"priority" json:"priority"`
	TeamID           int64                 `db:"team_id" json:"team_id"`
	TeamName         string                `db:"team_name" json:"team_name"`
	Area             string                `db:"area" json:"area"`
	RecordAlertDay   int                   `db:"record_alert_day" json:"record_alert_day"`
	LastRecordDate   string                `db:"last_record_date" json:"last_record_date"`
	Progress         int                   `db:"progress" json:"progress"`
	Status           int                   `db:"status" json:"status"`
	ActiveEvents     []ActiveEventResponse `json:"active_events"`
	Created          time.Time             `db:"created" json:"created"`
	CreatedBy        string                `db:"created_by" json:"created_by"`
	Updated          time.Time             `db:"updated" json:"updated"`
	UpdatedBy        string                `db:"updated_by" json:"updated_by"`
}

type ActiveEventResponse struct {
	EventID    int64              `db:"event_id" json:"event_id"`
	EventName  string             `db:"event_name" json:"event_name"`
	Status     int                `db:"status" json:"status"`
	ActiveType string             `json:"active_type"`
	AssignType int                `db:"assign_type" json:"assign_type"`
	AuditType  int                `db:"audit_type" json:"audit_type"`
	Actives    []AssignToResponse `json:"actives"`
}

type AssignToResponse struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
type ProjectReportNew struct {
	Name           string   `json:"name" binding:"required"`
	ReportDate     string   `json:"report_date" binding:"required,datetime=2006-01-02"`
	Content        string   `json:"content" binding:"required"`
	Links          []string `json:"links" binding:"omitempty"`
	User           string   `json:"user" swaggerignore:"true"`
	UserID         int64    `json:"user_id" swaggerignore:"true"`
	OrganizationID int64    `json:"organization_id" swaggerignore:"true"`
}

type ProjectReportFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	Status         string `form:"status" binding:"required,oneof=all active"`
	OrganizationID int64  `json:"organization_id" swaggerignore:"true"`
	UserID         int64  `json:"user_id" swaggerignore:"true"`
}

type ProjectReportResponse struct {
	ID           int64                             `db:"id" json:"id"`
	ProjectID    int64                             `db:"project_id" json:"project_id"`
	ProjectName  string                            `db:"project_name" json:"project_name"`
	UserID       int64                             `db:"user_id" json:"user_id"`
	Name         string                            `db:"name" json:"name"`
	ReportDate   string                            `db:"report_date" json:"report_date"`
	Content      string                            `db:"content" json:"content"`
	Status       int                               `db:"status" json:"status"`
	Updated      time.Time                         `db:"updated" json:"updated"`
	Links        []string                          `json:"links"`
	Avatar       string                            `db:"avatar" json:"avatar"`
	Username     string                            `db:"user_name" json:"user_name"`
	PositionName string                            `db:"position_name" json:"position_name"`
	Views        []ProjectReportMemberViewResponse `json:"views"`
	Viewed       bool                              `json:"viewed"`
}

type ProjectRecordNew struct {
	Name           string   `json:"name" binding:"required"`
	RecordDate     string   `json:"record_date" binding:"required,datetime=2006-01-02"`
	Content        string   `json:"content" binding:"required"`
	Plan           string   `json:"plan" binding:"omitempty"`
	Photos         []string `json:"photos" binding:"omitempty"`
	User           string   `json:"user" swaggerignore:"true"`
	UserID         int64    `json:"user_id" swaggerignore:"true"`
	OrganizationID int64    `json:"organization_id" swaggerignore:"true"`
}

type ProjectRecordFilter struct {
	PageId         int   `form:"page_id" binding:"required,min=1"`
	PageSize       int   `form:"page_size" binding:"required,min=5,max=200"`
	OrganizationID int64 `json:"organization_id" swaggerignore:"true"`
	UserID         int64 `json:"user_id" swaggerignore:"true"`
}

type ProjectRecordResponse struct {
	ID           int64     `db:"id" json:"id"`
	ProjectID    int64     `db:"project_id" json:"project_id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	Name         string    `db:"name" json:"name"`
	RecordDate   string    `db:"record_date" json:"record_date"`
	Content      string    `db:"content" json:"content"`
	Plan         string    `db:"plan" json:"plan"`
	Status       int       `db:"status" json:"status"`
	Updated      time.Time `db:"updated" json:"updated"`
	Photos       []string  `json:"photos"`
	Username     string    `db:"user_name" json:"user_name"`
	PositionName string    `db:"position_name" json:"position_name"`
	Avatar       string    `db:"avatar" json:"avatar"`
}
type ProjectReportViewResponse struct {
	ID              int64     `db:"id" json:"id"`
	ProjectID       int64     `db:"project_id" json:"project_id"`
	ProjectReportID int64     `db:"project_report_id" json:"project_report_id"`
	ViewerID        int64     `db:"viewer_id" json:"viewer_id"`
	ViewerName      string    `db:"viewer_name" json:"viewer_name"`
	Created         time.Time `db:"created" json:"created"`
}

type ProjectReportMemberViewResponse struct {
	UserID   int64     `json:"user_id"`
	UserName string    `json:"user_name"`
	Avatar   string    `json:"avatar"`
	Viewed   bool      `json:"viewed"`
	ViewTime time.Time `json:"view_time"`
}

type NodeAudit struct {
	AuditLevel int
	AuditTo    int64
}
