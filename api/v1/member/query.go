package member

import (
	"github.com/jmoiron/sqlx"
)

type memberQuery struct {
	conn *sqlx.DB
}

func NewMemberQuery(connection *sqlx.DB) *memberQuery {
	return &memberQuery{
		conn: connection,
	}
}

func (r *memberQuery) GetMembersByProjectID(projectID int64) (*[]MemberResponse, error) {
	var res []MemberResponse
	err := r.conn.Select(&res, `
		SELECT m.user_id, u.name 
		FROM project_members m 
		LEFT JOIN users u 
		ON m.user_id = u.id 
		WHERE m.project_id = ? 
		AND m.status = ? 
		`, projectID, 1)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
