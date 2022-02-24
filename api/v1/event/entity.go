package event

import "time"

type Event struct {
	ID        int64          `db:"id" json:"id"`
	ProjectID int64          `db:"project_id" json:"project_id"`
	Name      string         `db:"name" json:"name"`
	PreID     int64          `db:"pre_id" json:"pre_id"`
	Status    int            `db:"status" json:"status"`
	Assign    *[]EventAssign `json:"assign"`
	Created   time.Time      `db:"created" json:"created"`
	CreatedBy string         `db:"created_by" json:"created_by"`
	Updated   time.Time      `db:"updated" json:"updated"`
	UpdatedBy string         `db:"updated_by" json:"updated_by"`
}

type EventAssign struct {
	ID         int64     `db:"id" json:"id"`
	EventID    int64     `db:"event_id" json:"event_id"`
	AssignType int       `db:"assign_type" json:"assign_type"`
	AssignTo   string    `db:"assign_to" json:"assign_to"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}
