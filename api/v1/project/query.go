package project

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type projectQuery struct {
	conn *sqlx.DB
}

func NewProjectQuery(connection *sqlx.DB) ProjectQuery {
	return &projectQuery{
		conn: connection,
	}
}

type ProjectQuery interface {
	//Project Management
	GetProjectByID(int64, int64) (*Project, error)
	GetProjectCount(ProjectFilter, int64) (int, error)
	GetProjectList(ProjectFilter, int64) (*[]Project, error)
}

func (r *projectQuery) GetProjectByID(id int64, organizationID int64) (*Project, error) {
	var project Project
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&project, "SELECT * FROM projects WHERE id = ? AND organization_id = ?", id, organizationID)
	} else {
		err = r.conn.Get(&project, "SELECT * FROM projects WHERE id = ? ", id)
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectQuery) GetProjectCount(filter ProjectFilter, organizationID int64) (int, error) {
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
		FROM projects 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *projectQuery) GetProjectList(filter ProjectFilter, organizationID int64) (*[]Project, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var projects []Project
	err := r.conn.Select(&projects, `
		SELECT * 
		FROM projects 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}
