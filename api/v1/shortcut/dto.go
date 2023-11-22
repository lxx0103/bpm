package shortcut

type ShortcutFilter struct {
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	ShortcutType   int    `form:"shortcut_type" binding:"omitempty,min=1"`
	Keyword        string `form:"keyword" binding:"omitempty,max=64,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type ShortcutNew struct {
	OrganizationID   int64  `json:"organization_id" binding:"omitempty,min=1"`
	ShortcutType     int64  `json:"shortcut_type" binding:"required"`
	Content          string `json:"content" binding:"required"`
	ShortcutTypeName string `json:"shortcut_type_name" swaggerignore:"true"`
	User             string `json:"user" swaggerignore:"true"`
	// UserID           int64  `json:"user_id" swaggerignore:"true"`
}

type ShortcutUpdate struct {
	Content          string `json:"content" binding:"required"`
	User             string `json:"user" swaggerignore:"true"`
	ShortcutTypeName string `json:"shortcut_type_name" swaggerignore:"true"`
	UserID           int64  `json:"user_id" swaggerignore:"true"`
}
type ShortcutID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ShortcutResponse struct {
	ID               int64  `db:"id" json:"id"`
	OrganizationID   int64  `db:"organization_id" json:"organization_id"`
	OrganizationName string `db:"organization_name" json:"organization_name"`
	ShortcutType     int64  `db:"shortcut_type" json:"shortcut_type"`
	ShortcutTypeName string `db:"shortcut_type_name" json:"shortcut_type_name"`
	Content          string `db:"content" json:"content"`
	Status           int    `db:"status" json:"status"`
}

type ShortcutTypeFilter struct {
	OrganizationID int64 `form:"organization_id" binding:"required,min=1"`
	ParentID       int64 `form:"parent_id" binding:"omitempty"`
}

type ShortcutTypeNew struct {
	OrganizationID int64  `json:"organization_id" binding:"omitempty,min=1"`
	ParentID       int64  `json:"parent_id" binding:"omitempty"`
	Name           string `json:"name" binding:"required"`
	User           string `json:"user" swaggerignore:"true"`
	// UserID           int64  `json:"user_id" swaggerignore:"true"`
}

type ShortcutTypeUpdate struct {
	ParentID int64  `json:"parent_id" binding:"omitempty"`
	Name     string `json:"name" binding:"required"`
	User     string `json:"user" swaggerignore:"true"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}
type ShortcutTypeID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type ShortcutTypeResponse struct {
	ID               int64  `db:"id" json:"id"`
	OrganizationID   int64  `db:"organization_id" json:"organization_id"`
	OrganizationName string `db:"organization_name" json:"organization_name"`
	ParentID         int64  `db:"parent_id" json:"parent_id"`
	Name             string `db:"name" json:"name"`
	Status           int    `db:"status" json:"status"`
}
