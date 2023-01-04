package common

type BrandFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type BrandNew struct {
	Name string `json:"name" binding:"required,min=1,max=64"`
	User string `json:"user" swaggerignore:"true"`
}

type BrandID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type BrandResponse struct {
	ID     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Status int    `db:"status" json:"status"`
}

type MaterialFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type MaterialNew struct {
	Name string `json:"name" binding:"required,min=1,max=64"`
	User string `json:"user" swaggerignore:"true"`
}

type MaterialID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type MaterialResponse struct {
	ID     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Status int    `db:"status" json:"status"`
}
