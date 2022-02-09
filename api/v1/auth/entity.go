package auth

import "time"

type User struct {
	ID             int64     `db:"id" json:"id"`
	Type           int       `db:"type" json:"type"`
	Identifier     string    `db:"identifier" json:"identifier"`
	Credential     string    `db:"credential" json:"credential"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	PositionID     int64     `db:"position_id" json:"position_id"`
	RoleID         int64     `db:"role_id" json:"role_id"`
	Name           string    `db:"name" json:"name"`
	Email          string    `db:"email" json:"email"`
	Gender         string    `db:"gender" json:"gender"`
	Phone          string    `db:"phone" json:"phone"`
	Birthday       string    `db:"birthday" json:"birthday"`
	Address        string    `db:"address" json:"address"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
type Role struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Priority  int64     `db:"priority" json:"priority"`
	Status    string    `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
type UserAPI struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Method    string    `db:"method" json:"method"`
	Route     string    `db:"route" json:"route"`
	Enabled   string    `db:"enabled" json:"enabled"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
type UserMenu struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Action    string    `db:"action" json:"action"`
	Title     string    `db:"title" json:"title"`
	Path      string    `db:"path" json:"path"`
	Component string    `db:"component" json:"component"`
	IsHidden  int64     `db:"is_hidden" json:"is_hidden"`
	ParentID  int64     `db:"parent_id" json:"parent_id"`
	Enabled   int64     `db:"enabled" json:"enabled"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type UserMenuAPI struct {
	ID        int64     `db:"id" json:"id"`
	MenuID    int64     `db:"menu_id" json:"menu_id"`
	APIID     int64     `db:"api_id" json:"api_id"`
	Enabled   string    `db:"enabled" json:"enabled"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type UserRoleMenu struct {
	ID        int64     `db:"id" json:"id"`
	MenuID    int64     `db:"menu_id" json:"menu_id"`
	APIID     int64     `db:"api_id" json:"api_id"`
	Enabled   string    `db:"enabled" json:"enabled"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
