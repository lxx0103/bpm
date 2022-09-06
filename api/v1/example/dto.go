package example

type ExampleFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	Style          string `form:"style" binding:"omitempty"`
	Type           string `form:"type" binding:"omitempty"`
	Room           string `form:"room" binding:"omitempty"`
	ExampleType    int    `form:"example_type" binding:"omitempty"`
	Status         int    `form:"status" binding:"omitempty"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ExampleNew struct {
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Cover          string `json:"cover" binding:"omitempty"`
	OrganizationID int64  `json:"organization_id" binding:"omitempty,min=1"`
	Notes          string `json:"notes" binding:"omitempty"`
	Description    string `json:"description" binding:"omitempty"`
	Style          string `json:"style" binding:"omitempty"`
	Type           string `json:"type" binding:"omitempty"`
	Room           string `json:"room" binding:"omitempty"`
	ExampleType    int    `json:"example_type" binding:"omitempty"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	User           string `json:"user" swaggerignore:"true"`
}

type ExampleID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ExampleResponse struct {
	ID               int64  `db:"id" json:"id"`
	OrganizationID   int64  `db:"organization_id" json:"organization_id"`
	OrganizationName string `db:"organization_name" json:"organization_name"`
	ExampleType      int    `db:"example_type" json:"example_type"`
	Name             string `db:"name" json:"name"`
	Cover            string `db:"cover" json:"cover"`
	Notes            string `db:"notes" json:"notes"`
	Style            string `db:"style" json:"style"`
	Type             string `db:"type" json:"type"`
	Room             string `db:"room" json:"room"`
	Status           int    `db:"status" json:"status"`
	Description      string `db:"description" json:"description"`
}

type ExampleListResponse struct {
	ID               int64  `db:"id" json:"id"`
	OrganizationID   int64  `db:"organization_id" json:"organization_id"`
	OrganizationName string `db:"organization_name" json:"organization_name"`
	Name             string `db:"name" json:"name"`
	Cover            string `db:"cover" json:"cover"`
	Notes            string `db:"notes" json:"notes"`
	Style            string `db:"style" json:"style"`
	Type             string `db:"type" json:"type"`
	Room             string `db:"room" json:"room"`
	Status           int    `db:"status" json:"status"`
}
