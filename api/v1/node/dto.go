package node

type NodeFilter struct {
	Name       string `form:"name" binding:"omitempty,max=64,min=1"`
	TemplateID int64  `form:"template_id" binding:"required,min=1"`
	PageId     int    `form:"page_id" binding:"required,min=1"`
	PageSize   int    `form:"page_size" binding:"required,min=5,max=200"`
}

type NodeNew struct {
	TemplateID int64   `json:"template_id" binding:"required,min=1"`
	Name       string  `json:"name" binding:"required,min=1,max=64"`
	PreID      []int64 `json:"pre_id" binding:"required"`
	Assignable int     `json:"assignable" binding:"required,oneof=1 2"`
	AssignType int     `json:"assign_type" binding:"required,oneof=1 2 3"`
	AssignTo   []int64 `json:"assign_to" binding:"required"`
	User       string  `json:"user" swaggerignore:"true"`
}
type NodeUpdate struct {
	Name       string  `json:"name" binding:"omitempty,min=1,max=64"`
	PreID      []int64 `json:"pre_id" binding:"omitempty"`
	Assignable int     `json:"assignable" binding:"omitempty,oneof=1 2"`
	AssignType int     `json:"assign_type" binding:"omitempty,oneof=1 2 3"`
	AssignTo   []int64 `json:"assign_to" binding:"omitempty"`
	JsonData   string  `json:"json_data" binding:"required,json"`
	User       string  `json:"user" swaggerignore:"true"`
}

type NodeID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type MyNodeFilter struct {
	Status string `form:"status" binding:"required,oneof=all active"`
}

type MyNodeResponse struct {
}
