package component

import "time"

type Component struct {
	ID            int64     `db:"id" json:"id"`
	EventID       int64     `db:"event_id" json:"event_id"`
	Sort          int       `db:"sort" json:"sort"`
	ComponentType string    `db:"component_type" json:"component_type"`
	Name          string    `db:"name" json:"name"`
	Value         string    `db:"value" json:"value"`
	DefaultValue  string    `db:"default_value" json:"default_value"`
	Patterns      string    `db:"patterns" json:"patterns"`
	Required      int       `db:"required" json:"required"`
	Status        int       `db:"status" json:"status"`
	JsonData      string    `db:"json_data" json:"json_data"`
	Created       time.Time `db:"created" json:"created"`
	CreatedBy     string    `db:"created_by" json:"created_by"`
	Updated       time.Time `db:"updated" json:"updated"`
	UpdatedBy     string    `db:"updated_by" json:"updated_by"`
}
