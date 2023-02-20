package position

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type positionQuery struct {
	conn *sqlx.DB
}

func NewPositionQuery(connection *sqlx.DB) PositionQuery {
	return &positionQuery{
		conn: connection,
	}
}

type PositionQuery interface {
	//Position Management
	GetPositionByID(int64, int64) (*Position, error)
	GetPositionCount(PositionFilter) (int, error)
	GetPositionList(PositionFilter) (*[]PositionResponse, error)
}

func (r *positionQuery) GetPositionByID(id int64, organizationID int64) (*Position, error) {
	var position Position
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&position, "SELECT * FROM positions WHERE id = ? AND organization_id = ? AND status > 0", id, organizationID)
	} else {
		err = r.conn.Get(&position, "SELECT * FROM positions WHERE id = ? AND status > 0", id)
	}
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *positionQuery) GetPositionCount(filter PositionFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM positions 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *positionQuery) GetPositionList(filter PositionFilter) (*[]PositionResponse, error) {
	where, args := []string{"p.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "p.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var positions []PositionResponse
	err := r.conn.Select(&positions, `
		SELECT p.id as id, p.name as name, p.status as status, p.organization_id as organization_id, o.name as organization_name
		FROM positions p
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &positions, nil
}
