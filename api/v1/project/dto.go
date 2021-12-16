package project

type ProjectFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ProjectNew struct {
	OrganizationID int64  `json:"organization_id" binding:"required,min=1"`
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	User           string `json:"user" swaggerignore:"true"`
}

type ProjectID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
