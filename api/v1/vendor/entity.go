package vendor

import "time"

type Vendor struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Phone       string    `db:"phone" json:"phone"`
	Address     string    `db:"address" json:"address"`
	Longitude   string    `db:"longitude" json:"longitude"`
	Latitude    string    `db:"latitude" json:"latitude"`
	Description string    `db:"description" json:"description"`
	Cover       string    `db:"cover" json:"cover"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}
