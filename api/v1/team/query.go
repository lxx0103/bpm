package team

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type teamQuery struct {
	conn *sqlx.DB
}

func NewTeamQuery(connection *sqlx.DB) *teamQuery {
	return &teamQuery{
		conn: connection,
	}
}

func (r *teamQuery) GetTeamByID(id int64, organizationID int64) (*Team, error) {
	var team Team
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&team, "SELECT * FROM teams WHERE id = ? AND organization_id = ? AND status > 0", id, organizationID)
	} else {
		err = r.conn.Get(&team, "SELECT * FROM teams WHERE id = ? AND status > 0", id)
	}
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamQuery) GetTeamCount(filter TeamFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Status; v == "active" {
		where, args = append(where, "status = ?"), append(args, 1)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM teams 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *teamQuery) GetTeamList(filter TeamFilter) (*[]TeamResponse, error) {
	where, args := []string{"p.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "p.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	if v := filter.Status; v == "active" {
		where, args = append(where, "p.status = ?"), append(args, 1)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var teams []TeamResponse
	err := r.conn.Select(&teams, `
		SELECT p.id as id, p.name as name, p.leader as leader, p.phone as phone, p.status as status, p.organization_id as organization_id, o.name as organization_name
		FROM teams p
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &teams, nil
}
