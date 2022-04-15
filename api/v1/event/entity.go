package event

import "time"

type Event struct {
	ID           int64          `db:"id" json:"id"`
	ProjectID    int64          `db:"project_id" json:"project_id"`
	Name         string         `db:"name" json:"name"`
	Assignable   int            `db:"assignable" json:"assignable"`
	AssignType   int            `db:"assign_type" json:"assign_type"`
	NodeID       int64          `db:"node_id" json:"node_id"`
	PreID        *[]EventPre    `json:"pre_id"`
	NeedAudit    int            `db:"need_audit" json:"need_audit"`
	AuditType    int            `db:"audit_type" json:"audit_type"`
	Audit        *[]EventAudit  `json:"audit"`
	CompleteTime string         `db:"complete_time" json:"complete_time"`
	CompleteUser string         `db:"complete_user" json:"complete_user"`
	AuditTime    string         `db:"audit_time" json:"audit_time"`
	AuditUser    string         `db:"audit_user" json:"audit_user"`
	Status       int            `db:"status" json:"status"`
	Assign       *[]EventAssign `json:"assign"`
	Created      time.Time      `db:"created" json:"created"`
	CreatedBy    string         `db:"created_by" json:"created_by"`
	Updated      time.Time      `db:"updated" json:"updated"`
	UpdatedBy    string         `db:"updated_by" json:"updated_by"`
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
	ID        int64     `db:"id" json:"id"`
	EventID   int64     `db:"event_id" json:"event_id"`
	AuditType int       `db:"audit_type" json:"audit_type"`
	AuditTo   int64     `db:"audit_to" json:"audit_to"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
