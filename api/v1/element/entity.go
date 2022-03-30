package element

import "time"

type Element struct {
	ID           int64     `db:"id" json:"id"`
	NodeID       int64     `db:"node_id" json:"node_id"`
	Sort         int       `db:"sort" json:"sort"`
	ElementType  string    `db:"element_type" json:"element_type"`
	Name         string    `db:"name" json:"name"`
	Value        string    `db:"value" json:"value"`
	DefaultValue string    `db:"default_value" json:"default_value"`
	Patterns     string    `db:"patterns" json:"patterns"`
	Required     int       `db:"required" json:"required"`
	Status       int       `db:"status" json:"status"`
	JsonData     string    `db:"json_data" json:"json_data"`
	Created      time.Time `db:"created" json:"created"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	Updated      time.Time `db:"updated" json:"updated"`
	UpdatedBy    string    `db:"updated_by" json:"updated_by"`
}
