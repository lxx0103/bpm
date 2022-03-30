package element

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type elementQuery struct {
	conn *sqlx.DB
}

func NewElementQuery(connection *sqlx.DB) ElementQuery {
	return &elementQuery{
		conn: connection,
	}
}

type ElementQuery interface {
	//Element Management
	GetElementByID(id int64) (*Element, error)
	GetElementCount(filter ElementFilter) (int, error)
	GetElementList(filter ElementFilter) (*[]Element, error)
}

func (r *elementQuery) GetElementByID(id int64) (*Element, error) {
	var element Element
	err := r.conn.Get(&element, "SELECT * FROM elements WHERE status > 0 AND id = ? ", id)
	if err != nil {
		return nil, err
	}
	return &element, nil
}

func (r *elementQuery) GetElementCount(filter ElementFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.NodeID; v != 0 {
		where, args = append(where, "node_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM elements 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *elementQuery) GetElementList(filter ElementFilter) (*[]Element, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.NodeID; v != 0 {
		where, args = append(where, "node_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var elements []Element
	err := r.conn.Select(&elements, `
		SELECT * 
		FROM elements 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &elements, nil
}
