package example

import "time"

type Example struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Cover          string    `db:"cover" json:"cover"`
	Notes          string    `db:"notes" json:"notes"`
	Description    string    `db:"description" json:"description"`
	Style          string    `db:"style" json:"style"`
	Type           string    `db:"type" json:"type"`
	Room           string    `db:"room" json:"room"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
