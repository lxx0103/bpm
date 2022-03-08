package member

import "time"

type Member struct {
	ID        int64           `db:"id" json:"id"`
	ProjectID int64           `db:"project_id" json:"project_id"`
	Name      string          `db:"name" json:"name"`
	PreID     *[]MemberPre    `json:"pre_id"`
	Status    int             `db:"status" json:"status"`
	Assign    *[]MemberAssign `json:"assign"`
	Created   time.Time       `db:"created" json:"created"`
	CreatedBy string          `db:"created_by" json:"created_by"`
	Updated   time.Time       `db:"updated" json:"updated"`
	UpdatedBy string          `db:"updated_by" json:"updated_by"`
}

type MemberAssign struct {
	ID         int64     `db:"id" json:"id"`
	MemberID   int64     `db:"member_id" json:"member_id"`
	AssignType int       `db:"assign_type" json:"assign_type"`
	AssignTo   string    `db:"assign_to" json:"assign_to"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}

type MemberPre struct {
	ID        int64     `db:"id" json:"id"`
	MemberID  int64     `db:"member_id" json:"member_id"`
	PreID     int64     `db:"pre_id" json:"pre_id"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
