package project

type ProjectFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	Type           int    `form:"type" binding:"omitempty,oneof=1 2"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ProjectNew struct {
	Name       string  `json:"name" binding:"required,min=1,max=64"`
	TemplateID int64   `json:"template_id" binding:"required,min=1"`
	ClientID   int64   `json:"client_id" binding:"omitempty,min=1"`
	Type       int     `json:"type" swaggerignore:"true"`
	Location   string  `json:"location" binding:"omitempty,min=1,max=64"`
	Longitude  float64 `json:"longitude" binding:"omitempty"`
	Latitude   float64 `json:"latitude" binding:"omitempty"`
	User       string  `json:"user" swaggerignore:"true"`
	UserID     int64   `json:"user_id" swaggerignore:"true"`
}

type ProjectID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ProjectUpdate struct {
	Name      string  `json:"name" binding:"omitempty,min=1,max=64"`
	ClientID  int64   `json:"client_id" binding:"omitempty,min=1"`
	Location  string  `form:"location" binding:"omitempty,min=1,max=64"`
	Longitude float64 `form:"longitude" binding:"omitempty"`
	Latitude  float64 `form:"latitude" binding:"omitempty"`
	User      string  `json:"user" swaggerignore:"true"`
	UserID    int64   `json:"user_id" swaggerignore:"true"`
}
type MyProjectFilter struct {
	Status   string `form:"status" binding:"required,oneof=all active"`
	Type     int    `form:"type" binding:"omitempty,oneof=1 2"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}
type AssignedProjectFilter struct {
	// Status   string `form:"status" binding:"required,oneof=all active"`
	Type     int `form:"type" binding:"omitempty,oneof=1 2"`
	PageId   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=5,max=200"`
}
