package event

import "time"

type Event struct {
	ID              int64          `db:"id" json:"id"`
	ProjectID       int64          `db:"project_id" json:"project_id"`
	Name            string         `db:"name" json:"name"`
	Assignable      int            `db:"assignable" json:"assignable"`
	AssignType      int            `db:"assign_type" json:"assign_type"`
	NodeID          int64          `db:"node_id" json:"node_id"`
	PreID           *[]EventPre    `json:"pre_id"`
	NeedAudit       int            `db:"need_audit" json:"need_audit"`
	AuditLevel      int            `db:"audit_level" json:"audit_level"`
	AuditType       int            `db:"audit_type" json:"audit_type"`
	Audit           *[]EventAudit  `json:"audit"`
	CompleteTime    string         `db:"complete_time" json:"complete_time"`
	CompleteUser    string         `db:"complete_user" json:"complete_user"`
	AuditTime       string         `db:"audit_time" json:"audit_time"`
	AuditContent    string         `db:"audit_content" json:"audit_content"`
	AuditUser       string         `db:"audit_user" json:"audit_user"`
	NeedCheckin     int            `db:"need_checkin" json:"need_checkin"`
	CheckinDistance int            `db:"checkin_distance" json:"checkin_distance"`
	Sort            int            `db:"sort" json:"sort"`
	CanReview       int            `db:"can_review" json:"can_review"`
	Deadline        string         `db:"deadline" json:"deadline"`
	Status          int            `db:"status" json:"status"`
	Assign          *[]EventAssign `json:"assign"`
	AuditFile       []string       `json:"audit_file"`
	Created         time.Time      `db:"created" json:"created"`
	CreatedBy       string         `db:"created_by" json:"created_by"`
	Updated         time.Time      `db:"updated" json:"updated"`
	UpdatedBy       string         `db:"updated_by" json:"updated_by"`
}

type EventAssign struct {
	ID         int64     `db:"id" json:"id"`
	EventID    int64     `db:"event_id" json:"event_id"`
	AssignType int       `db:"assign_type" json:"assign_type"`
	AssignTo   int64     `db:"assign_to" json:"assign_to"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}

type EventPre struct {
	ID        int64     `db:"id" json:"id"`
	EventID   int64     `db:"event_id" json:"event_id"`
	PreID     int64     `db:"pre_id" json:"pre_id"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
type EventAudit struct {
	ID         int64     `db:"id" json:"id"`
	EventID    int64     `db:"event_id" json:"event_id"`
	AuditLevel int       `db:"audit_level" json:"audit_level"`
	AuditType  int       `db:"audit_type" json:"audit_type"`
	AuditTo    int64     `db:"audit_to" json:"audit_to"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}
type EventCheckin struct {
	ID          int64     `db:"id" json:"id"`
	EventID     int64     `db:"event_id" json:"event_id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	UserName    string    `db:"user_name" json:"user_name"`
	CheckinType int       `db:"checkin_type" json:"checkin_type"`
	CheckinTime time.Time `db:"checkin_time" json:"checkin_time"`
	Distance    int       `db:"distance" json:"distance"`
	Longitude   float64   `db:"longitude" json:"longitude"`
	Latitude    float64   `db:"latitude" json:"latitude"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

type EventAuditHistory struct {
	ID           int64     `db:"id" json:"id"`
	EventID      int64     `db:"event_id" json:"event_id"`
	AuditTime    string    `db:"audit_time" json:"audit_time"`
	AuditContent string    `db:"audit_content" json:"audit_content"`
	AuditUser    string    `db:"audit_user" json:"audit_user"`
	Status       int       `db:"status" json:"status"`
	Created      time.Time `db:"created" json:"created"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	Updated      time.Time `db:"updated" json:"updated"`
	UpdatedBy    string    `db:"updated_by" json:"updated_by"`
}

type EventReview struct {
	ID            int64     `db:"id" json:"id"`
	EventID       int64     `db:"event_id" json:"event_id"`
	Result        int       `db:"result" json:"result"`
	Content       string    `db:"content" json:"content"`
	Link          string    `db:"link" json:"link"`
	Status        int       `db:"status" json:"status"`
	HandleTime    string    `db:"handle_time" json:"handle_time"`
	HandleContent string    `db:"handle_content" json:"handle_content"`
	HandleUser    string    `db:"handle_user" json:"handle_user"`
	Created       time.Time `db:"created" json:"created"`
	CreatedBy     string    `db:"created_by" json:"created_by"`
	Updated       time.Time `db:"updated" json:"updated"`
	UpdatedBy     string    `db:"updated_by" json:"updated_by"`
}

type EventAuditFile struct {
	ID        int64     `db:"id" json:"id"`
	EventID   int64     `db:"event_id" json:"event_id"`
	Link      string    `db:"link" json:"link"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type EventHistoryFile struct {
	ID        int64     `db:"id" json:"id"`
	HistoryID int64     `db:"history_id" json:"history_id"`
	Link      string    `db:"link" json:"link"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
