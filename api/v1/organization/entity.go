package organization

import "time"

type Organization struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Logo        string    `db:"logo" json:"logo"`
	Logo2       string    `db:"logo2" json:"logo2"`
	Description string    `db:"description" json:"description"`
	Phone       string    `db:"phone" json:"phone"`
	Contact     string    `db:"contact" json:"contact"`
	Address     string    `db:"address" json:"address"`
	City        string    `db:"city" json:"city"`
	Type        int       `db:"type" json:"type"`
	UserLimit   int       `db:"user_limit" json:"user_limit"`
	ExpiryDate  string    `db:"expiry_date" json:"expiry_date"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}
