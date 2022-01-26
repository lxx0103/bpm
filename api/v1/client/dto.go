package client

type ClientFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ClientNew struct {
	Name    string `json:"name" binding:"required,min=1,max=64"`
	Phone   string `json:"phone" binding:"required,min=6,max=64"`
	Address string `json:"address" binding:"omitempty,max=255"`
	Status  int    `json:"status" binding:"required,oneof=1 2"`
	User    string `json:"user" swaggerignore:"true"`
}

type ClientID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
