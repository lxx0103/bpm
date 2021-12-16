package component

type ComponentFilter struct {
	EventID  int64  `form:"event_id" binding:"omitempty,min=1"`
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ComponentNew struct {
	EventID     int64  `json:"event_id" binding:"required,min=1"`
	Type        string `json:"type" binding:"required,min=1,max=32"`
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Description string `json:"description" binding:"required,min=1,max=255"`
	Info        string `json:"info" binding:"required,min=1"`
	Status      int    `json:"status" binding:"required,oneof=1 2"`
	User        string `json:"user" swaggerignore:"true"`
}

type ComponentID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
