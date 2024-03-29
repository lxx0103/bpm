package node

import "time"

type Node struct {
	ID          int64         `db:"id" json:"id"`
	TemplateID  int64         `db:"template_id" json:"template_id"`
	Name        string        `db:"name" json:"name"`
	Assignable  int           `db:"assignable" json:"assignable"`
	AssignType  int           `db:"assign_type" json:"assign_type"`
	Assign      *[]NodeAssign `json:"assign"`
	NeedAudit   int           `db:"need_audit" json:"need_audit"`
	AuditType   int           `db:"audit_type" json:"audit_type"`
	Audit       *[]NodeAudit  `json:"audit"`
	JsonData    string        `db:"json_data" json:"json_data"`
	PreID       *[]NodePre    `json:"pre_id"`
	NeedCheckin int           `db:"need_checkin" json:"need_checkin"`
	Sort        int           `db:"sort" json:"sort"`
	CanReview   int           `db:"can_review" json:"can_review"`
	Status      int           `db:"status" json:"status"`
	Created     time.Time     `db:"created" json:"created"`
	CreatedBy   string        `db:"created_by" json:"created_by"`
	Updated     time.Time     `db:"updated" json:"updated"`
	UpdatedBy   string        `db:"updated_by" json:"updated_by"`
}

type NodeAssign struct {
	ID         int64     `db:"id" json:"id"`
	NodeID     int64     `db:"node_id" json:"node_id"`
	AssignType int       `db:"assign_type" json:"assign_type"`
	AssignTo   int64     `db:"assign_to" json:"assign_to"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}

type NodePre struct {
	ID        int64     `db:"id" json:"id"`
	NodeID    int64     `db:"node_id" json:"node_id"`
	PreID     int64     `db:"pre_id" json:"pre_id"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type NodeAudit struct {
	ID         int64     `db:"id" json:"id"`
	NodeID     int64     `db:"node_id" json:"node_id"`
	AuditLevel int       `db:"audit_level" json:"audit_level"`
	AuditType  int       `db:"audit_type" json:"audit_type"`
	AuditTo    int64     `db:"audit_to" json:"audit_to"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}
