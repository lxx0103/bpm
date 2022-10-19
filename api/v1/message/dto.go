package message

type MessageFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type MessageNew struct {
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Priority int    `json:"priority" binding:"required,min=1"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	User     string `json:"user" swaggerignore:"true"`
}

type MessageID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
