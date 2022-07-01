package client

import "time"

type Client struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	Phone          string    `db:"phone" json:"phone"`
	Address        string    `db:"address" json:"address"`
	Name           string    `db:"name" json:"name"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
