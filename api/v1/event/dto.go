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
	Sort        int    `json:"sort" binding:"required, min=1"`
	CanReview   int    `json:"can_review" binding:"required,oneof=1 2"`
	NodeID      int64  `json:"node_id" binding:"required,min=1"`
	User        string `json:"user" swaggerignore:"true"`
}
type EventUpdate struct {
	AssignType int     `json:"assign_type" binding:"omitempty,oneof=1 2"`
	AssignTo   []int64 `json:"assign_to" binding:"omitempty"`
	NeedAudit  int     `json:"need_audit" binding:"omitempty,oneof=1 2"`
	AuditType  int     `json:"audit_type" binding:"omitempty,oneof=1 2"`
	AuditTo    []int64 `json:"audit_to" binding:"omitempty"`
	AuditMore  []struct {
		AuditLevel int     `json:"audit_level" binding:"required,min=2"`
		AuditType  int     `json:"audit_type" binding:"required,oneof=1 2"`
		AuditTo    []int64 `json:"audit_to" binding:"required"`
	} `json:"audit_more" binding:"omitempty"`
	User string `json:"user" swaggerignore:"true"`
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
	ID           int64                 `db:"id" json:"id"`
	ProjectID    int64                 `db:"project_id" json:"project_id"`
	ProjectName  string                `db:"project_name" json:"project_name"`
	Name         string                `db:"name" json:"name"`
	CompleteTime string                `db:"complete_time" json:"complete_time"`
	CompleteUser string                `db:"complete_user" json:"complete_user"`
	AuditTime    string                `db:"audit_time" json:"audit_time"`
	AuditUser    string                `db:"audit_user" json:"audit_user"`
	AuditContent string                `db:"audit_content" json:"audit_content"`
	NeedCheckin  int                   `db:"need_checkin" json:"need_checkin"`
	Sort         int                   `db:"sort" json:"sort"`
	Status       int                   `db:"status" json:"status"`
	Priority     int                   `db:"priority" json:"priority"`
	Deadline     string                `db:"deadline" json:"deadline"`
	CanReview    int                   `db:"can_review" json:"can_review"`
	IsActive     int                   `db:"is_active" json:"is_active"`
	NeedAudit    int                   `db:"need_audit" json:"need_audit"`
	AuditLevel   int                   `db:"audit_level" json:"audit_level"`
	AuditType    int                   `db:"audit_type" json:"audit_type"`
	Audit        *[][]AssignToResponse `json:"audit"`
	Assignable   int                   `db:"assignable" json:"assignable"`
	AssignType   int                   `db:"assign_type" json:"assign_type"`
	Assign       *[]AssignToResponse   `json:"assign"`
}

type AssignToResponse struct {
	ID         int64  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	AuditType  string `db:"audit_type" json:"audit_type"`
	AudtiLevel string `db:"audit_level" json:"audit_level"`
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
	Content    string `json:"content" binding:"omitempty,max=255"`
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
	From           string `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To             string `form:"to" binding:"omitempty,datetime=2006-01-02"`
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
type EventAuditHistoryResponse struct {
	ID           int64  `db:"id" json:"id"`
	EventID      int64  `db:"event_id" json:"event_id"`
	AuditTime    string `db:"audit_time" json:"audit_time"`
	AuditContent string `db:"audit_content" json:"audit_content"`
	AuditUser    string `db:"audit_user" json:"audit_user"`
	Status       int    `db:"status" json:"status"`
}

type EventReviewNew struct {
	Result     int    `json:"result" binding:"required,oneof=1 2"`
	Content    string `json:"content" binding:"omitempty,max=255"`
	Link       string `json:"link" binding:"omitempty"`
	User       string `json:"user" swaggerignore:"true"`
	UserID     int64  `json:"user_id" swaggerignore:"true"`
	PositionID int64  `json:"position_id" swaggerignore:"true"`
}

type EventReviewResponse struct {
	ID            int64  `db:"id" json:"id"`
	EventID       int64  `db:"event_id" json:"event_id"`
	Result        string `db:"result" json:"result"`
	Content       string `db:"content" json:"content"`
	Link          string `db:"link" json:"link"`
	Status        int    `db:"status" json:"status"`
	Created       string `db:"created" json:"created"`
	HandleTime    string `db:"handle_time" json:"handle_time"`
	HandleContent string `db:"handle_content" json:"handle_content"`
	HandleUser    string `db:"handle_user" json:"handle_user"`
}

type EventDeadlineNew struct {
	Deadline string `json:"deadline" binding:"omitempty,datetime=2006-01-02"`
	User     string `json:"user" swaggerignore:"true"`
}

type HandleReviewInfo struct {
	Result     int    `json:"result" binding:"required,oneof=2 3"`
	Content    string `json:"content" binding:"omitempty,max=255"`
	User       string `json:"user" swaggerignore:"true"`
	UserID     int64  `json:"user_id" swaggerignore:"true"`
	PositionID int64  `json:"position_id" swaggerignore:"true"`
}

type NodeAudit struct {
	AuditLevel int
	AuditTo    []int64
}
