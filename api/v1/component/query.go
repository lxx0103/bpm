package component

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type componentQuery struct {
	conn *sqlx.DB
}

func NewComponentQuery(connection *sqlx.DB) ComponentQuery {
	return &componentQuery{
		conn: connection,
	}
}

type ComponentQuery interface {
	//Component Management
	GetComponentByID(id int64) (*Component, error)
	GetComponentCount(filter ComponentFilter) (int, error)
	GetComponentList(filter ComponentFilter) (*[]Component, error)
}

func (r *componentQuery) GetComponentByID(id int64) (*Component, error) {
	var component Component
	err := r.conn.Get(&component, "SELECT * FROM event_components WHERE id = ? ", id)
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *componentQuery) GetComponentCount(filter ComponentFilter) (int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.EventID; v != 0 {
		where, args = append(where, "event_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM event_components 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *componentQuery) GetComponentList(filter ComponentFilter) (*[]Component, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.EventID; v != 0 {
		where, args = append(where, "event_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var components []Component
	err := r.conn.Select(&components, `
		SELECT * 
		FROM event_components 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &components, nil
}
