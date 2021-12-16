package request

type PageInfo struct {
	PageSize int `json:"page_size" form:"page_size" binding:"required,oneof=5 10 15 20"`
	PageID   int `json:"page_id" form:"page_id" binding:"required,min=1"`
}
