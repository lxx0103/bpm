package auth

type WechatCredential struct {
	OpenID     string `json:"openid" binding:"required"`
	SessionKey string `json:"session_key" binding:"required"`
	UnionID    string `json:"union_id"`
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}
type SigninRequest struct {
	AuthType       int    `json:"auth_type" binding:"required,oneof=1 2 3"`
	Identifier     string `json:"identifier" binding:"required"`
	Credential     string `json:"credential" binding:"omitempty,min=6"`
	OrganizationID int64  `json:"organization_id" binding:"omitempty,min=1"`
}
type SigninResponse struct {
	Token string `json:"token"`
	User  UserResponse
}

type SignupRequest struct {
	OrganizationID int64  `json:"organization_id" binding:"required,min=1"`
	Identifier     string `json:"identifier" binding:"required"`
	Credential     string `json:"credential" binding:"required,min=6"`
}

type RoleFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type RoleNew struct {
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Priority int    `json:"priority" binding:"required,min=1"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	User     string `json:"user" swaggerignore:"true"`
}

type RoleID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UserFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	Type           string `form:"type" binding:"omitempty,oneof=wx admin"`
	OrganizationID int64  `form:"organization_id" binding:"omitempty,min=1"`
	PageId         int    `form:"page_id" binding:"required,min=1"`
	PageSize       int    `form:"page_size" binding:"required,min=5,max=200"`
}

type APIFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	Route    string `form:"route" binding:"omitempty,max=128,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type APINew struct {
	Name   string `json:"name" binding:"required,min=1,max=64"`
	Route  string `json:"route" binding:"required,min=1,max=128"`
	Method string `json:"method" binding:"required,oneof=post put get"`
	Status int    `json:"status" binding:"required,oneof=1 2"`
	User   string `json:"user" swaggerignore:"true"`
}

type APIID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type MenuFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	OnlyTop  bool   `form:"only_top" binding:"omitempty"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type MenuNew struct {
	Name      string `json:"name" binding:"required,min=1,max=64"`
	Action    string `json:"action" binding:"omitempty,min=1,max=64"`
	Title     string `json:"title" binding:"required,min=1,max=64"`
	Path      string `json:"path" binding:"omitempty,min=1,max=128"`
	Component string `json:"component" binding:"omitempty,min=1,max=255"`
	IsHidden  int64  `json:"is_hidden" binding:"required,oneof=1 2"`
	ParentID  int64  `json:"parent_id" binding:"required,min=-1"`
	Status    int    `json:"status" binding:"required,oneof=1 2"`
	User      string `json:"user" swaggerignore:"true"`
}

type MenuUpdate struct {
	Name      string `json:"name" binding:"omitempty,min=1,max=64"`
	Action    string `json:"action" binding:"omitempty,min=1,max=64"`
	Title     string `json:"title" binding:"omitempty,min=1,max=64"`
	Path      string `json:"path" binding:"omitempty,min=1,max=128"`
	Component string `json:"component" binding:"omitempty,min=1,max=255"`
	IsHidden  int64  `json:"is_hidden" binding:"omitempty,oneof=1 2"`
	ParentID  int64  `json:"parent_id" binding:"omitempty,min=-1"`
	Status    int    `json:"status" binding:"required,oneof=1 2"`
	User      string `json:"user" swaggerignore:"true"`
}

type MenuID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type RoleMenu struct {
	IDS []int64 `json:"ids" binding:"required"`
}
type RoleMenuNew struct {
	IDS  []int64 `json:"ids" binding:"required"`
	User string  `json:"_" swaggerignore:"true"`
}

type MenuAPIID struct {
	IDS []int64 `json:"ids" binding:"required"`
}

type MenuAPINew struct {
	IDS  []int64 `json:"ids" binding:"required"`
	User string  `json:"_" swaggerignore:"true"`
}

type MyMenuDetail struct {
	Name      string         `json:"name" binding:"required,min=1,max=64"`
	Action    string         `json:"action" binding:"omitempty,min=1,max=64"`
	Title     string         `json:"title" binding:"required,min=1,max=64"`
	Path      string         `json:"path" binding:"omitempty,min=1,max=128"`
	Component string         `json:"component" binding:"omitempty,min=1,max=255"`
	IsHidden  int64          `json:"is_hidden" binding:"required,oneof=1 2"`
	Status    int            `json:"status" binding:"required,oneof=1 2"`
	Items     []MyMenuDetail `json:"items"`
}

type UserID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type UserUpdate struct {
	RoleID     int64  `json:"role_id" binding:"omitempty,min=1"`
	PositionID int64  `json:"position_id" binding:"omitempty,min=1"`
	Name       string `json:"name" binding:"omitempty,min=2"`
	Email      string `json:"email" binding:"omitempty,email"`
	Gender     string `json:"gender" binding:"omitempty,min=1"`
	Phone      string `json:"phone" binding:"omitempty,min=1"`
	Birthday   string `json:"birthday" binding:"omitempty,datetime=2006-01-02"`
	Address    string `json:"address" binding:"omitempty,min=1"`
	Avatar     string `json:"avatar" binding:"omitempty"`
	Status     int    `json:"status" binding:"omitempty,min=1"`
	User       string `json:"user" swaggerignore:"true"`
}
type UserResponse struct {
	ID               int64  `db:"id" json:"id"`
	Type             int    `db:"type" json:"type"`
	Identifier       string `db:"identifier" json:"identifier"`
	OrganizationID   int64  `db:"organization_id" json:"organization_id"`
	OrganizationName string `db:"organization_name" json:"organization_name"`
	PositionID       int64  `db:"position_id" json:"position_id"`
	RoleID           int64  `db:"role_id" json:"role_id"`
	Name             string `db:"name" json:"name"`
	Email            string `db:"email" json:"email"`
	Gender           string `db:"gender" json:"gender"`
	Phone            string `db:"phone" json:"phone"`
	Birthday         string `db:"birthday" json:"birthday"`
	Address          string `db:"address" json:"address"`
	Avatar           string `db:"avatar" json:"avatar"`
	Status           int    `db:"status" json:"status"`
}

type PasswordUpdate struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	User        string `json:"user" swaggerignore:"true"`
	UserID      int64  `json:"user_id" swaggerignore:"true"`
}

type WxmoduleFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type WxmoduleNew struct {
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Code     string `json:"code" binding:"omitempty,min=1,max=64"`
	ParentID int64  `json:"parent_id" binding:"required,min=-1"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	User     string `json:"user" swaggerignore:"true"`
}

type WxmoduleUpdate struct {
	Name     string `json:"name" binding:"omitempty,min=1,max=64"`
	Code     string `json:"code" binding:"omitempty,min=1,max=64"`
	ParentID int64  `json:"parent_id" binding:"omitempty,min=-1"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	User     string `json:"user" swaggerignore:"true"`
}

type WxmoduleID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type PositionWxmodule struct {
	IDS []int64 `json:"ids" binding:"required"`
}
type PositionWxmoduleNew struct {
	IDS  []int64 `json:"ids" binding:"required"`
	User string  `json:"_" swaggerignore:"true"`
}

type MyWxmoduleDetail struct {
	Name   string             `json:"name" binding:"required,min=1,max=64"`
	Code   string             `json:"code" binding:"omitempty,min=1,max=64"`
	Status int                `json:"status" binding:"required,oneof=1 2"`
	Items  []MyWxmoduleDetail `json:"items"`
}
