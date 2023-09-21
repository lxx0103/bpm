package vendors

type VendorsFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	Brand    string `form:"brand" binding:"omitempty,max=64,min=1"`
	Material string `form:"material" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type VendorsNew struct {
	Name        string          `json:"name" binding:"required,min=1,max=64"`
	Material    []int64         `json:"material"`
	Brand       []int64         `json:"brand"`
	Qrcode      []VendorsQrcode `json:"qrcode"`
	Contact     string          `json:"contact" binding:"omitempty"`
	Phone       string          `json:"phone" binding:"omitempty,max=64"`
	Address     string          `json:"address" binding:"omitempty,max=255"`
	Longitude   float64         `json:"longitude" binding:"omitempty"`
	Latitude    float64         `json:"latitude" binding:"omitempty"`
	Cover       string          `json:"cover" binding:"omitempty,max=255"`
	Picture     []string        `json:"picture"`
	Description string          `json:"description" binding:"omitempty"`
	User        string          `json:"user" swaggerignore:"true"`
}

type VendorsQrcode struct {
	Type string `db:"type" json:"type" binding:"required"`
	Name string `db:"name" json:"name" binding:"required"`
}

type VendorsID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type VendorsResponse struct {
	ID          int64             `db:"id" json:"id"`
	Name        string            `db:"name" json:"name"`
	Material    []VendorsMaterial `json:"material"`
	Brand       []VendorsBrand    `json:"brand"`
	Contact     string            `db:"contact" json:"contact"`
	Phone       string            `db:"phone" json:"phone"`
	Address     string            `db:"address" json:"address"`
	Longitude   string            `db:"longitude" json:"longitude"`
	Latitude    string            `db:"latitude" json:"latitude"`
	Cover       string            `db:"cover" json:"cover"`
	Description string            `db:"description" json:"description"`
	Picture     []string          `json:"picture"`
	Qrcode      []VendorsQrcode   `json:"qrcode"`
	Status      int               `db:"status" json:"status"`
}

type VendorsMaterial struct {
	MaterialID   int64  `db:"material_id" json:"material_id"`
	MaterialName string `db:"material_name" json:"material_name"`
}

type VendorsBrand struct {
	BrandID   int64  `db:"brand_id" json:"brand_id"`
	BrandName string `db:"brand_name" json:"brand_name"`
}
