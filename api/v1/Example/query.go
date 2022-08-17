package example

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type exampleQuery struct {
	conn *sqlx.DB
}

func NewExampleQuery(connection *sqlx.DB) ExampleQuery {
	return &exampleQuery{
		conn: connection,
	}
}

type ExampleQuery interface {
	//Example Management
	GetExampleByID(int64, int64) (*Example, error)
	GetExampleCount(ExampleFilter) (int, error)
	GetExampleList(ExampleFilter) (*[]ExampleListResponse, error)
}

func (r *exampleQuery) GetExampleByID(id int64, organizationID int64) (*Example, error) {
	var example Example
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&example, "SELECT * FROM examples WHERE id = ? AND organization_id = ?", id, organizationID)
	} else {
		err = r.conn.Get(&example, "SELECT * FROM examples WHERE id = ? ", id)
	}
	if err != nil {
		return nil, err
	}
	return &example, nil
}

func (r *exampleQuery) GetExampleCount(filter ExampleFilter) (int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Style; v != "" {
		where, args = append(where, "style = ?"), append(args, v)
	}
	if v := filter.Type; v != "" {
		where, args = append(where, "type = ?"), append(args, v)
	}
	if v := filter.Room; v != "" {
		where, args = append(where, "room = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM examples 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *exampleQuery) GetExampleList(filter ExampleFilter) (*[]ExampleListResponse, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "e.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "e.organization_id = ?"), append(args, v)
	}
	if v := filter.Style; v != "" {
		where, args = append(where, "e.style = ?"), append(args, v)
	}
	if v := filter.Type; v != "" {
		where, args = append(where, "e.type = ?"), append(args, v)
	}
	if v := filter.Room; v != "" {
		where, args = append(where, "e.room = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var examples []ExampleListResponse
	err := r.conn.Select(&examples, `
		SELECT e.id, e.name, e.cover, e.style, e.type, e.room, e.notes, e.status, e.organization_id, o.name as organization_name
		FROM examples e
		LEFT JOIN organizations o
		ON e.organization_id = o.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &examples, nil
}
