package event

type EventFilter struct {
	Name      string `form:"name" binding:"omitempty,max=64,min=1"`
	ProjectID int64  `form:"project_id" binding:"omitempty,min=1"`
	PageId    int    `form:"page_id" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=5,max=200"`
}

type EventNew struct {
	ProjectID  int64   `json:"project_id" binding:"required,min=1"`
	Name       string  `json:"name" binding:"required,min=1,max=64"`
	PreID      int64   `json:"pre_id" binding:"required,min=-1"`
	AssignType int     `json:"assign_type" binding:"required,oneof=1 2"`
	AssignTo   []int64 `json:"assign_to" binding:"required"`
	Status     int     `json:"status" binding:"required,oneof=1 2 3"`
	User       string  `json:"user" swaggerignore:"true"`
}
type EventUpdate struct {
	Name       string  `json:"name" binding:"omitempty,min=1,max=64"`
	PreID      int64   `json:"pre_id" binding:"omitempty,min=1"`
	AssignType int     `json:"assign_type" binding:"omitempty,oneof=1 2"`
	AssignTo   []int64 `json:"assign_to" binding:"omitempty"`
	Status     int     `json:"status" binding:"omitempty,oneof=1 2 3"`
	User       string  `json:"user" swaggerignore:"true"`
}

type EventID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
