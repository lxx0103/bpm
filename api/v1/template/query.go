package template

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type templateQuery struct {
	conn *sqlx.DB
}

func NewTemplateQuery(connection *sqlx.DB) TemplateQuery {
	return &templateQuery{
		conn: connection,
	}
}

type TemplateQuery interface {
	//Template Management
	GetTemplateByID(int64, int64) (*Template, error)
	GetTemplateCount(TemplateFilter, int64) (int, error)
	GetTemplateList(TemplateFilter, int64) (*[]Template, error)
}

func (r *templateQuery) GetTemplateByID(id int64, organizationID int64) (*Template, error) {
	var template Template
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&template, "SELECT * FROM templates WHERE status > 0 AND id = ? AND organization_id = ?", id, organizationID)
	} else {
		err = r.conn.Get(&template, "SELECT * FROM templates WHERE status > 0 AND id = ? ", id)
	}
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateQuery) GetTemplateCount(filter TemplateFilter, organizationID int64) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	} else if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM templates 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *templateQuery) GetTemplateList(filter TemplateFilter, organizationID int64) (*[]Template, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	} else if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var templates []Template
	err := r.conn.Select(&templates, `
		SELECT * 
		FROM templates 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &templates, nil
}
