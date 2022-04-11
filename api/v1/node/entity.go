package node

import "time"

type Node struct {
	ID         int64         `db:"id" json:"id"`
	TemplateID int64         `db:"template_id" json:"template_id"`
	Name       string        `db:"name" json:"name"`
	Assignable int           `db:"assignable" json:"assignable"`
	AssignType int           `db:"assign_type" json:"assign_type"`
	JsonData   string        `db:"json_data" json:"json_data"`
	PreID      *[]NodePre    `json:"pre_id"`
	Status     int           `db:"status" json:"status"`
	Assign     *[]NodeAssign `json:"assign"`
	Created    time.Time     `db:"created" json:"created"`
	CreatedBy  string        `db:"created_by" json:"created_by"`
	Updated    time.Time     `db:"updated" json:"updated"`
	UpdatedBy  string        `db:"updated_by" json:"updated_by"`
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