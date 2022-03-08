package member

type MemberFilter struct {
	ProjectID int64 `form:"project_id" binding:"required,min=1"`
}

type MemberNew struct {
	ProjectID int64   `json:"project_id" binding:"required,min=1"`
	UserID    []int64 `json:"user_id" binding:"required"`
	User      string  `json:"user" swaggerignore:"true"`
}

type MemberResponse struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
}
