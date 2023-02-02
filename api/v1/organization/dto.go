package organization

type OrganizationFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	City     string `form:"city" binding:"omitempty"`
	Type     int    `form:"type" binding:"omitempty,oneof=1 2"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type OrganizationNew struct {
	Name        string               `json:"name" binding:"required,min=1,max=64"`
	Description string               `json:"description" binding:"required"`
	Logo        string               `json:"logo" binding:"omitempty"`
	Contact     string               `json:"contact" binding:"omitempty"`
	Phone       string               `json:"phone" binding:"omitempty"`
	Address     string               `json:"address" binding:"omitempty"`
	City        string               `json:"city" binding:"required"`
	Type        int                  `json:"type" binding:"required,oneof=1 2"`
	Status      int                  `json:"status" binding:"required,oneof=1 2"`
	Qrcode      []OrganizationQrcode `json:"qrcode"`
	User        string               `json:"user" swaggerignore:"true"`
}

type OrganizationResponse struct {
	ID          int64                `db:"id" json:"id"`
	Name        string               `db:"name" json:"name"`
	Logo        string               `db:"logo" json:"logo"`
	Description string               `db:"description" json:"description"`
	Phone       string               `db:"phone" json:"phone"`
	Contact     string               `db:"contact" json:"contact"`
	Address     string               `db:"address" json:"address"`
	City        string               `db:"city" json:"city"`
	Type        int                  `db:"type" json:"type"`
	Status      int                  `db:"status" json:"status"`
	Qrcode      []OrganizationQrcode `json:"qrcode"`
}
type OrganizationQrcode struct {
	Type string `db:"type" json:"type" binding:"required"`
	Name string `db:"name" json:"name" binding:"required"`
}
type OrganizationID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type QrcodeFilter struct {
	Path   string `json:"path" binding:"required,max=128,min=1"`
	Source string `json:"source" binding:"required,oneof=bpm portal"`
}

type WechatToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type OrganizationExampleResponse struct {
	ID          int64                `db:"id" json:"id"`
	Name        string               `db:"name" json:"name"`
	Logo        string               `db:"logo" json:"logo"`
	Description string               `db:"description" json:"description"`
	Phone       string               `db:"phone" json:"phone"`
	Contact     string               `db:"contact" json:"contact"`
	Address     string               `db:"address" json:"address"`
	City        string               `db:"city" json:"city"`
	Type        int                  `db:"type" json:"type"`
	Status      int                  `db:"status" json:"status"`
	Examples    []ExampleResponse    `json:"examples"`
	Qrcode      []OrganizationQrcode `json:"qrcode"`
}
type ExampleResponse struct {
	ID     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Cover  string `db:"cover" json:"cover"`
	Status int    `db:"status" json:"status"`
}
