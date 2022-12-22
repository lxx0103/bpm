package organization

import "time"

type Organization struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Logo        string    `db:"logo" json:"logo"`
	Description string    `db:"description" json:"description"`
	Phone       string    `db:"phone" json:"phone"`
	Contact     string    `db:"contact" json:"contact"`
	Address     string    `db:"address" json:"address"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}
