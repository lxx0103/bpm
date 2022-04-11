package upload

type UploadFilter struct {
	OrganizationID int64  `form:"organization_id" binding:"required,min=1"`
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}
