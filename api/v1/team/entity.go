package team

import "time"

type Team struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Leader         string    `db:"leader" json:"leader"`
	Phone          string    `db:"phone" json:"phone"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
