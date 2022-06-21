package organization

type OrganizationFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type OrganizationNew struct {
	Name   string `json:"name" binding:"required,min=1,max=64"`
	Status int    `json:"status" binding:"required,oneof=1 2"`
	User   string `json:"user" swaggerignore:"true"`
}

type OrganizationID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type QrcodeFilter struct {
	Path string `json:"path" binding:"required,max=128,min=1"`
}

type WechatToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}
