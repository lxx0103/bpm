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
	GetPositionCount(PositionFilter, int64) (int, error)
	GetPositionList(PositionFilter, int64) (*[]Position, error)
}

func (r *positionQuery) GetPositionByID(id int64, organizationID int64) (*Position, error) {
	var position Position
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&position, "SELECT * FROM positions WHERE id = ? AND organization_id = ?", id, organizationID)
	} else {
		err = r.conn.Get(&position, "SELECT * FROM positions WHERE id = ? ", id)
	}
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *positionQuery) GetPositionCount(filter PositionFilter, organizationID int64) (int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
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

func (r *positionQuery) GetPositionList(filter PositionFilter, organizationID int64) (*[]Position, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var positions []Position
	err := r.conn.Select(&positions, `
		SELECT * 
		FROM positions 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &positions, nil
}
