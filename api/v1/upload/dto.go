package upload

type UploadFilter struct {
	OrganizationID int64  `form:"organization_id" binding:"required,min=1"`
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type KeyFilter struct {
	APPID  string `form:"app_id" binding:"required,min=1"`
	Bucket string `form:"bucket" binding:"required,min=1"`
}

type KeyRes struct {
	TmpSecretId  string `json:"TmpSecretId"`
	TmpSecretKey string `json:"TmpSecretKey"`
	Token        string `json:"Token"`
	StartTime    int    `json:"StartTime"`
	ExpiredTime  int    `json:"ExpiredTime"`
	Expiration   string `json:"Expiration"`
}
