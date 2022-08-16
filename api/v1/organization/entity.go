package organization

import "time"

type Organization struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}
