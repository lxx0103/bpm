package component

import "time"

type Component struct {
	ID          int64     `db:"id" json:"id"`
	EventID     int64     `db:"event_id" json:"event_id"`
	Type        string    `db:"type" json:"type"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Info        string    `db:"info" json:"info"`
	Value       string    `db:"value" json:"value"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}
