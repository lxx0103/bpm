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
	//WX
	GetProjectListByCreate(string, int64, MyProjectFilter) (*[]Project, error)
	GetProjectCountByCreate(string, int64, MyProjectFilter) (int, error)
	GetProjectListByAssigned(AssignedProjectFilter, int64, int64, int64) (*[]Project, error)
	GetProjectCountByAssigned(AssignedProjectFilter, int64, int64, int64) (int, error)
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
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	} else if v := filter.OrganizationID; v != 0 {
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
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	} else if v := filter.OrganizationID; v != 0 {
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

func (r *projectQuery) GetProjectListByCreate(userName string, organization_id int64, filter MyProjectFilter) (*[]Project, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if filter.Status == "all" {
		where, args = append(where, "status > ?"), append(args, 0)
	} else {
		where, args = append(where, "status = ?"), append(args, 1)
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	where, args = append(where, "created_by = ?"), append(args, userName)
	where, args = append(where, "organization_id = ?"), append(args, organization_id)
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

func (r *projectQuery) GetProjectCountByCreate(userName string, organization_id int64, filter MyProjectFilter) (int, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if filter.Status == "all" {
		where, args = append(where, "status > ?"), append(args, 0)
	} else {
		where, args = append(where, "status = ?"), append(args, 1)
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	where, args = append(where, "created_by = ?"), append(args, userName)
	where, args = append(where, "organization_id = ?"), append(args, organization_id)
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

func (r *projectQuery) GetProjectListByAssigned(filter AssignedProjectFilter, userID int64, positionID int64, organizationID int64) (*[]Project, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	var projects []Project
	err := r.conn.Select(&projects, `
		SELECT * FROM projects WHERE id IN 
		(
			SELECT project_id FROM events WHERE id IN 
			(
				SELECT event_id FROM event_assigns WHERE ((assign_type = 1 and assign_to = ?) or (assign_type = 2 and assign_to = ?)) AND status > 0
			) AND status > 0
		)
		AND `+strings.Join(where, " AND ")+`
		AND status > 0 AND organization_id = ? 
		ORDER BY ID DESC
		LIMIT ?, ?
	`, positionID, userID, args, organizationID, filter.PageId*filter.PageSize-filter.PageSize, filter.PageSize)
	return &projects, err
}

func (r *projectQuery) GetProjectCountByAssigned(filter AssignedProjectFilter, userID int64, positionID int64, organizationID int64) (int, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) FROM projects WHERE id IN 
		(
			SELECT project_id FROM events WHERE id IN 
			(
				SELECT event_id FROM event_assigns WHERE ((assign_type = 1 and assign_to = ?) or (assign_type = 2 and assign_to = ?)) AND status > 0
			) AND status > 0
		)
		AND `+strings.Join(where, " AND ")+`
		AND status > 0 AND organization_id = ? 
	`, positionID, userID, args, organizationID)
	return count, err
}
