package event

type EventFilter struct {
	Name      string `form:"name" binding:"omitempty,max=64,min=1"`
	ProjectID int64  `form:"project_id" binding:"omitempty,min=1"`
	PageId    int    `form:"page_id" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=5,max=200"`
}

type EventNew struct {
	ProjectID   int64  `json:"project_id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Assignable  int    `json:"assignable" binding:"required,oneof=1 2"`
	AssignType  int    `json:"assign_type" binding:"required,oneof=1 2 3"`
	NeedAudit   int    `json:"need_audit" binding:"required,oneof=1 2"`
	AuditType   int    `json:"audit_type" binding:"required,oneof=1 2"`
	NeedCheckin int    `json:"need_checkin" binding:"required,oneof=1 2"`
	NodeID      int64  `json:"node_id" binding:"required,min=1"`
	User        string `json:"user" swaggerignore:"true"`
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
	AuditTime    string `db:"audit_time" json:"audit_time"`
	AuditUser    string `db:"audit_user" json:"audit_user"`
	AuditContent string `db:"audit_content" json:"audit_content"`
	NeedCheckin  int    `db:"need_checkin" json:"need_checkin"`
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

type AuditEventInfo struct {
	Result     int    `json:"result" binding:"required,oneof=1 2"`
	Content    string `json:"content" binding:"required,max=255"`
	User       string `json:"user" swaggerignore:"true"`
	UserID     int64  `json:"user_id" swaggerignore:"true"`
	PositionID int64  `json:"position_id" swaggerignore:"true"`
}
type AssignedAuditFilter struct {
	Status string `form:"status" binding:"required,oneof=all active"`
}

type NewCheckin struct {
	Longitude      float64 `json:"longitude" binding:"required"`
	Latitude       float64 `json:"latitude" binding:"required"`
	User           string  `json:"user" swaggerignore:"true"`
	OrganizationID int64   `json:"organization_id" swaggerignore:"true"`
	PositionID     int64   `json:"position_id" swaggerignore:"true"`
	UserID         int64   `json:"user_id" swaggerignore:"true"`
	CheckinType    int     `json:"checkin_type" swaggerignore:"true"`
	Distance       int     `json:"distance" swaggerignore:"true"`
}

type CheckinFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	ProjectID      int64  `form:"project_id" binding:"omitempty,min=1"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	EventID        int64  `form:"event_id" binding:"omitempty,min=1"`
	UserID         int64  `form:"user_id" binding:"omitempty,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type CheckinResponse struct {
	Name             string  `db:"name" json:"name"`
	ProjectID        int64   `db:"project_id" json:"project_id"`
	ProjectName      string  `db:"project_name" json:"project_name"`
	EventID          int64   `db:"event_id" json:"event_id"`
	EventName        string  `db:"event_name" json:"event_name"`
	OrganizationID   int64   `db:"organization_id" json:"organization_id"`
	OrganizationName string  `db:"organization_name" json:"organization_name"`
	CheckinType      int     `db:"checkin_type" json:"checkin_type"`
	CheckinTime      string  `db:"checkin_time" json:"checkin_time"`
	Longitude        float64 `db:"longitude" json:"longitude"`
	Latitude         float64 `db:"latitude" json:"latitude"`
	Distance         int     `db:"distance" json:"distance"`
}
