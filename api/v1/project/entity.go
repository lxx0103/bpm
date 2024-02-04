package project

import "time"

type Project struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	TemplateID      int64     `db:"template_id" json:"template_id"`
	ClientID        int64     `db:"client_id" json:"client_id"`
	Name            string    `db:"name" json:"name"`
	Type            int       `db:"type" json:"type"`
	Location        string    `db:"location" json:"location"`
	Longitude       float64   `db:"longitude" json:"longitude"`
	Latitude        float64   `db:"latitude" json:"latitude"`
	CheckinDistance int       `db:"checkin_distance" json:"checkin_distance"`
	Priority        int       `db:"priority" json:"priority"`
	Progress        int       `db:"progress" json:"progress"`
	TeamID          int64     `db:"team_id" json:"team_id"`
	Area            string    `db:"area" json:"area"`
	RecordAlertDay  int       `db:"record_alert_day" json:"record_alert_day"`
	LastRecordDate  string    `db:"last_record_date" json:"last_record_date"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}

type ProjectReport struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	ProjectID      int64     `db:"project_id" json:"project_id"`
	ClientID       int64     `db:"client_id" json:"client_id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	ReportDate     string    `db:"report_date" json:"report_date"`
	Name           string    `db:"name" json:"name"`
	Content        string    `db:"content" json:"content"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type ProjectReportLink struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	ProjectID       int64     `db:"project_id" json:"project_id"`
	ProjectReportID int64     `db:"project_report_id" json:"project_report_id"`
	Link            string    `db:"link" json:"link"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}

type ProjectReportView struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	ProjectID       int64     `db:"project_id" json:"project_id"`
	ProjectReportID int64     `db:"project_report_id" json:"project_report_id"`
	ViewerID        int64     `db:"viewer_id" json:"viewer_id"`
	ViewerName      string    `db:"viewer_name" json:"viewer_name"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}

type ProjectRecord struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	ProjectID      int64     `db:"project_id" json:"project_id"`
	ClientID       int64     `db:"client_id" json:"client_id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	RecordDate     string    `db:"record_date" json:"record_date"`
	Name           string    `db:"name" json:"name"`
	Content        string    `db:"content" json:"content"`
	Plan           string    `db:"plan" json:"plan"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type ProjectRecordPhoto struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	ProjectID       int64     `db:"project_id" json:"project_id"`
	ProjectRecordID int64     `db:"project_record_id" json:"project_record_id"`
	Link            string    `db:"link" json:"link"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}
