package upload

import "time"

type Upload struct {
	ID             int64     `db:"id" json:"id"`
	Path           string    `db:"path" json:"path"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
