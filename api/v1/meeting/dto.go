package meeting

import "time"

type MeetingFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type MeetingNew struct {
	Name           string `json:"name" binding:"required,min=1,max=64"`
	OrganizationID int64  `json:"organization_id" binding:"omitempty,min=1"`
	Date           string `json:"date" binding:"required,datetime=2006-01-02"`
	Content        string `json:"content" binding:"required"`
	File           string `json:"file" binding:"omitempty"`
	User           string `json:"user" swaggerignore:"true"`
}

type MeetingID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type MeetingResponse struct {
	ID               int64     `db:"id" json:"id"`
	OrganizationID   int64     `db:"organization_id" json:"organization_id"`
	OrganizationName string    `db:"organization_name" json:"organization_name"`
	Date             time.Time `db:"date" json:"date"`
	Name             string    `db:"name" json:"name"`
	Content          string    `db:"content" json:"content"`
	File             string    `db:"file" json:"file"`
	Status           int       `db:"status" json:"status"`
}
