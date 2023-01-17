package example

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type exampleQuery struct {
	conn *sqlx.DB
}

func NewExampleQuery(connection *sqlx.DB) *exampleQuery {
	return &exampleQuery{
		conn: connection,
	}
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
	if v := filter.ExampleType; v != 0 {
		where, args = append(where, "example_type = ?"), append(args, v)
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
	if v := filter.Status; v != 0 {
		where, args = append(where, "status = ?"), append(args, v)
	}
	if v := filter.Mixed; v != "" {
		where = append(where, "(name like ? or building like ?)")
		args = append(args, "%"+v+"%")
		args = append(args, "%"+v+"%")
	}
	if v := filter.Priority; v == "index" {
		where, args = append(where, "priority > ?"), append(args, 0)
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
	if v := filter.ExampleType; v != 0 {
		where, args = append(where, "example_type = ?"), append(args, v)
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
	if v := filter.Status; v != 0 {
		where, args = append(where, "e.status = ?"), append(args, v)
	}
	if v := filter.Mixed; v != "" {
		where = append(where, "(e.name like ? or e.building like ?)")
		args = append(args, "%"+v+"%")
		args = append(args, "%"+v+"%")
	}
	if v := filter.Priority; v == "index" {
		where, args = append(where, "e.priority > ?"), append(args, 0)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var examples []ExampleListResponse
	err := r.conn.Select(&examples, `
		SELECT e.id, e.name, e.cover, e.style, e.type, e.room, e.notes, e.building, e.priority, e.status, e.organization_id, o.name as organization_name
		FROM examples e
		LEFT JOIN organizations o
		ON e.organization_id = o.id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY e.priority DESC, e.id DESC
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &examples, nil
}

func (r *exampleQuery) GetExampleMaterialList(exampleID int64) (*[]ExampleMaterialResponse, error) {
	var examples []ExampleMaterialResponse
	err := r.conn.Select(&examples, `
		SELECT em.id, 
		em.example_id, IFNULL(e.name, "") as example_name,
		em.material_id, IFNULL(m.name, "") as material_name,
		em.vendor_id, IFNULL(v.name, "") as vendor_name,
		em.brand_id, IFNULL(b.name, "") as brand_name,
		em.status
		FROM example_materials em 
		LEFT JOIN examples e ON em.example_id = e.id 
		LEFT JOIN materials m ON em.material_id = m.id 
		LEFT JOIN vendors v ON em.vendor_id = v.id 
		LEFT JOIN brands b ON em.brand_id = b.id 
		WHERE em.example_id = ?
		AND em.status > 0
	`, exampleID)
	return &examples, err
}

func (r *exampleQuery) GetExampleMaterialByID(exampleID, ID int64) (*ExampleMaterialResponse, error) {
	var examples ExampleMaterialResponse
	err := r.conn.Get(&examples, `
		SELECT em.id, 
		em.example_id, IFNULL(e.name, "") as example_name,
		em.material_id, IFNULL(m.name, "") as material_name,
		em.vendor_id, IFNULL(v.name, "") as vendor_name,
		em.brand_id, IFNULL(b.name, "") as brand_name,
		em.status
		FROM example_materials em 
		LEFT JOIN examples e ON em.example_id = e.id 
		LEFT JOIN materials m ON em.material_id = m.id 
		LEFT JOIN vendors v ON em.vendor_id = v.id 
		LEFT JOIN brands b ON em.brand_id = b.id 
		WHERE em.example_id = ?
		AND em.id = ?
		AND em.status > 0
	`, exampleID, ID)
	return &examples, err
}
