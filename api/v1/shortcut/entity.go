package shortcut

import "time"

type Shortcut struct {
	ID               int64     `db:"id" json:"id"`
	OrganizationID   int64     `db:"organization_id" json:"organization_id"`
	ShortcutType     int       `db:"shortcut_type" json:"shortcut_type"`
	ShortcutTypeName string    `db:"shortcut_type_name" json:"shortcut_type_name"`
	Content          string    `db:"content" json:"content"`
	Status           int       `db:"status" json:"status"`
	Created          time.Time `db:"created" json:"created"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	Updated          time.Time `db:"updated" json:"updated"`
	UpdatedBy        string    `db:"updated_by" json:"updated_by"`
}

type ShortcutType struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	ParentID       int64     `db:"parent_id" json:"parent_id"`
	Name           string    `db:"name" json:"name"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
