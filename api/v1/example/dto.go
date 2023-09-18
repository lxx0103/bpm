package example

type ExampleFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	Style          string `form:"style" binding:"omitempty"`
	Type           string `form:"type" binding:"omitempty"`
	Room           string `form:"room" binding:"omitempty"`
	ExampleType    int    `form:"example_type" binding:"omitempty"`
	Status         int    `form:"status" binding:"omitempty"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	Mixed          string `form:"mixed" binding:"omitempty"`
	Priority       string `form:"priority" binding:"omitempty,oneof=all index"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ExampleNew struct {
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Cover          string `json:"cover" binding:"omitempty"`
	OrganizationID int64  `json:"organization_id" binding:"omitempty,min=1"`
	Notes          string `json:"notes" binding:"omitempty"`
	Description    string `json:"description" binding:"omitempty"`
	Description2   string `json:"description2" binding:"omitempty"`
	Style          string `json:"style" binding:"omitempty"`
	Type           string `json:"type" binding:"omitempty"`
	Room           string `json:"room" binding:"omitempty"`
	FinderUserName string `json:"finder_user_name" binding:"omitempty"`
	FeedID         string `json:"feed_id" binding:"omitempty"`
	ExampleType    int    `json:"example_type" binding:"omitempty"`
	Priority       int    `json:"priority" binding:"omitempty"`
	Building       string `json:"building" binding:"omitempty"`
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
	FinderUserName   string `db:"finder_user_name" json:"finder_user_name"`
	FeedID           string `db:"feed_id" json:"feed_id"`
	Priority         int    `db:"priority" json:"priority"`
	Building         string `db:"building" json:"building"`
	Status           int    `db:"status" json:"status"`
	Description      string `db:"description" json:"description"`
	Description2     string `db:"description2" json:"description2"`
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
	Priority         int    `db:"priority" json:"priority"`
	Building         string `db:"building" json:"building"`
	Status           int    `db:"status" json:"status"`
}

type ExampleMaterialResponse struct {
	ID           int64  `db:"id" json:"id"`
	ExampleID    int64  `db:"example_id" json:"example_id"`
	ExampleName  string `db:"example_name" json:"example_name"`
	MaterialID   int64  `db:"material_id" json:"material_id"`
	MaterialName string `db:"material_name" json:"material_name"`
	VendorID     int64  `db:"vendor_id" json:"vendor_id"`
	VendorName   string `db:"vendor_name" json:"vendor_name"`
	BrandID      int64  `db:"brand_id" json:"brand_id"`
	BrandName    string `db:"brand_name" json:"brand_name"`
	Status       int    `db:"status" json:"status"`
}

type ExampleMaterialID struct {
	ID         int64 `uri:"id" binding:"required,min=1"`
	MaterialID int64 `uri:"material_id" binding:"required,min=1"`
}

type ExampleMaterialNew struct {
	MaterialID int64  `json:"material_id" binding:"required,min=1"`
	VendorID   int64  `json:"vendor_id" binding:"omitempty,min=1"`
	BrandID    int64  `json:"brand_id" binding:"omitempty,min=1"`
	User       string `json:"user" swaggerignore:"true"`
}
