package vendor

type VendorFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	Brand    string `form:"brand" binding:"omitempty,max=64,min=1"`
	Material string `form:"material" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type VendorNew struct {
	Name        string   `json:"name" binding:"required,min=1,max=64"`
	Material    []int64  `json:"material"`
	Brand       []int64  `json:"brand"`
	Contact     string   `json:"contact" binding:"omitempty"`
	Phone       string   `json:"phone" binding:"omitempty,max=64"`
	Address     string   `json:"address" binding:"omitempty,max=255"`
	Longitude   float64  `json:"longitude" binding:"omitempty"`
	Latitude    float64  `json:"latitude" binding:"omitempty"`
	Cover       string   `json:"cover" binding:"omitempty,max=64"`
	Picture     []string `json:"picture"`
	Description string   `json:"description" binding:"omitempty"`
	User        string   `json:"user" swaggerignore:"true"`
}

type VendorID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type VendorResponse struct {
	ID          int64            `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Material    []VendorMaterial `json:"material"`
	Brand       []VendorBrand    `json:"brand"`
	Contact     string           `db:"contact" json:"contact"`
	Phone       string           `db:"phone" json:"phone"`
	Address     string           `db:"address" json:"address"`
	Longitude   string           `db:"longitude" json:"longitude"`
	Latitude    string           `db:"latitude" json:"latitude"`
	Cover       string           `db:"cover" json:"cover"`
	Description string           `db:"description" json:"description"`
	Picture     []string         `json:"picture"`
	Status      int              `db:"status" json:"status"`
}

type VendorMaterial struct {
	MaterialID   int64  `db:"material_id" json:"material_id"`
	MaterialName string `db:"material_name" json:"material_name"`
}

type VendorBrand struct {
	BrandID   int64  `db:"brand_id" json:"brand_id"`
	BrandName string `db:"brand_name" json:"brand_name"`
}
