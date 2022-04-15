package event

type EventFilter struct {
	Name      string `form:"name" binding:"omitempty,max=64,min=1"`
	ProjectID int64  `form:"project_id" binding:"omitempty,min=1"`
	PageId    int    `form:"page_id" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=5,max=200"`
}

type EventNew struct {
	ProjectID  int64  `json:"project_id" binding:"required,min=1"`
	Name       string `json:"name" binding:"required,min=1,max=64"`
	Assignable int    `json:"assignable" binding:"required,oneof=1 2"`
	AssignType int    `json:"assign_type" binding:"required,oneof=1 2 3"`
	NeedAudit  int    `json:"need_audit" binding:"required,oneof=1 2"`
	AuditType  int    `json:"audit_type" binding:"required,oneof=1 2"`
	NodeID     int64  `json:"node_id" binding:"required,min=1"`
	User       string `json:"user" swaggerignore:"true"`
}
type EventUpdate struct {
	AssignType int     `json:"assign_type" binding:"omitempty,oneof=1 2"`
	AssignTo   []int64 `json:"assign_to" binding:"omitempty"`
	NeedAudit  int     `json:"need_audit" binding:"omitempty,oneof=1 2"`
	AuditType  int     `json:"audit_type" binding:"omitempty,oneof=1 2"`
	AuditTo    []int64 `json:"audit_to" binding:"omitempty"`
	User       string  `json:"user" swaggerignore:"true"`
}

type EventID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type AssignedEventFilter struct {
	Status string `form:"status" binding:"required,oneof=all active"`
}

type MyEventFilter struct {
	Status    string `form:"status" binding:"required,oneof=all active"`
	ProjectID int64  `form:"project_id" binding:"required,min=1"`
}

type MyEventResponse struct {
}

type MyEvent struct {
	ID           int64  `db:"id" json:"id"`
	ProjectID    int64  `db:"project_id" json:"project_id"`
	ProjectName  string `db:"project_name" json:"project_name"`
	Name         string `db:"name" json:"name"`
	CompleteTime string `db:"complete_time" json:"complete_time"`
	CompleteUser string `db:"complete_user" json:"complete_user"`
	Status       int    `db:"status" json:"status"`
}

type SaveEventInfo struct {
	Components []ComponentInfo `json:"component_info" binding:"required"`
	User       string          `json:"user" swaggerignore:"true"`
	UserID     int64           `json:"user_id" swaggerignore:"true"`
	PositionID int64           `json:"position_id" swaggerignore:"true"`
}

type ComponentInfo struct {
	ID    int64  `json:"id" binding:"required,min=1"`
	Value string `json:"value" binding:"required"`
}
