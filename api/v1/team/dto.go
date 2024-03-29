package team

type TeamFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	Status         string `form:"status" binding:"omitempty,oneof=active all"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type TeamNew struct {
	Name           string `json:"name" binding:"required,min=1,max=64"`
	OrganizationID int64  `json:"organization_id" binding:"omitempty,min=1"`
	Leader         string `json:"leader" binding:"required,min=1,max=64"`
	Phone          string `json:"phone" binding:"omitempty,min=1,max=64"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	User           string `json:"user" swaggerignore:"true"`
}

type TeamID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type TeamResponse struct {
	ID               int64  `db:"id" json:"id"`
	OrganizationID   int64  `db:"organization_id" json:"organization_id"`
	OrganizationName string `db:"organization_name" json:"organization_name"`
	Name             string `db:"name" json:"name"`
	Leader           string `db:"leader" json:"leader"`
	Phone            string `db:"phone" json:"phone"`
	Status           int    `db:"status" json:"status"`
}
