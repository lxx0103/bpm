package project

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type projectQuery struct {
	conn *sqlx.DB
}

func NewProjectQuery(connection *sqlx.DB) *projectQuery {
	return &projectQuery{
		conn: connection,
	}
}

func (r *projectQuery) GetProjectByID(id int64, organizationID int64) (*Project, error) {
	var project Project
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&project, "SELECT * FROM projects WHERE id = ? AND organization_id = ? AND status > 0", id, organizationID)
	} else {
		err = r.conn.Get(&project, "SELECT * FROM projects WHERE id = ? AND status > 0", id)
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

func (r *projectQuery) GetProjectList(filter ProjectFilter, organizationID int64) (*[]ProjectResponse, error) {
	where, args := []string{"p.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "p.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "p.type = ?"), append(args, v)
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	} else if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var projects []ProjectResponse
	err := r.conn.Select(&projects, `
		SELECT p.id as id, p.organization_id as organization_id, o.name as organization_name, p.client_id as client_id, IFNULL(c.name, "内部流程") as client_name, p.name as name, p.type as type, p.location as location, p.longitude as longitude, p.latitude as latitude, p.checkin_distance as checkin_distance, p.priority, p.status as status
		FROM projects p
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		LEFT JOIN clients c
		ON p.client_id = c.id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY p.id desc
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
		ORDER BY id desc
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
	// args = append(args, positionID)
	args = append(args, userID)
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	if v := filter.Status; v != 0 {
		where, args = append(where, "status = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	args = append(args, organizationID)
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var projects []Project
	// err := r.conn.Select(&projects, `
	// 	SELECT * FROM projects WHERE id IN
	// 	(
	// 		SELECT project_id FROM events WHERE id IN
	// 		(
	// 			SELECT event_id FROM event_assigns WHERE ((assign_type = 1 and assign_to = ?) or (assign_type = 2 and assign_to = ?)) AND status > 0
	// 		) AND status > 0
	// 	)
	// 	AND `+strings.Join(where, " AND ")+`
	// 	AND status > 0 AND organization_id = ?
	// 	ORDER BY ID DESC
	// 	LIMIT ?, ?
	// `, args...)
	err := r.conn.Select(&projects, `
		SELECT * FROM projects WHERE id IN 
		(
			SELECT project_id FROM project_members WHERE user_id = ? AND status > 0
		)
		AND `+strings.Join(where, " AND ")+`
		AND status > 0 AND organization_id = ? 
		ORDER BY ID DESC
		LIMIT ?, ?
	`, args...)
	return &projects, err
}

func (r *projectQuery) GetProjectCountByAssigned(filter AssignedProjectFilter, userID int64, positionID int64, organizationID int64) (int, error) {
	where, args := []string{"1=1"}, []interface{}{}
	// args = append(args, positionID)
	args = append(args, userID)
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	if v := filter.Status; v != 0 {
		where, args = append(where, "status = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	args = append(args, organizationID)
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) FROM projects WHERE id IN 
		(
			SELECT project_id FROM project_members WHERE user_id = ? AND status > 0
		)
		AND `+strings.Join(where, " AND ")+`
		AND status > 0 AND organization_id = ? 
	`, args...)
	return count, err
}

func (r *projectQuery) GetProjectListByClientID(userID int64, organization_id int64, filter MyProjectFilter) (*[]Project, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if filter.Status == "all" {
		where, args = append(where, "p.status > ?"), append(args, 0)
	} else {
		where, args = append(where, "p.status = ?"), append(args, 1)
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "p.type = ?"), append(args, v)
	}
	where, args = append(where, "c.user_id = ?"), append(args, userID)
	where, args = append(where, "p.organization_id = ?"), append(args, organization_id)
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var projects []Project
	err := r.conn.Select(&projects, `
		SELECT p.* 
		FROM projects p
		LEFT JOIN clients c
		ON p.client_id = c.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}

func (r *projectQuery) GetProjectCountByClientID(userID int64, organization_id int64, filter MyProjectFilter) (int, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if filter.Status == "all" {
		where, args = append(where, "p.status > ?"), append(args, 0)
	} else {
		where, args = append(where, "p.status = ?"), append(args, 1)
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "p.type = ?"), append(args, v)
	}
	where, args = append(where, "c.user_id = ?"), append(args, userID)
	where, args = append(where, "p.organization_id = ?"), append(args, organization_id)
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM projects p
		LEFT JOIN clients c
		ON p.client_id = c.id
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *projectQuery) GetProjectReportList(projectID int64, filter ProjectReportFilter) (*[]ProjectReportResponse, error) {
	where, args := []string{"pr.status > 0"}, []interface{}{}
	where, args = append(where, "pr.project_id = ?"), append(args, projectID)
	if filter.Status == "active" {
		where, args = append(where, "pr.status < ?"), append(args, 3)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "pr.name like ?"), append(args, "%"+v+"%")
	}
	var projectReports []ProjectReportResponse
	err := r.conn.Select(&projectReports, `
		SELECT pr.id, pr.user_id, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id DESC
	`, args...)
	return &projectReports, err
}

func (r *projectQuery) GetProjectReportByID(id int64, organizationID int64) (*ProjectReportResponse, error) {
	var report ProjectReportResponse
	if organizationID == 0 {
		err := r.conn.Get(&report, `
		SELECT pr.id, pr.project_id, pr.user_id, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id 
		WHERE pr.id = ? AND pr.status > 0 limit 1
		`, id)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.conn.Get(&report, `
		SELECT pr.id, pr.project_id, pr.user_id, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id 
		WHERE pr.id = ? AND pr.organization_id = ? AND pr.status > 0 limit 1`, id, organizationID)
		if err != nil {
			return nil, err
		}
	}
	return &report, nil
}

func (r *projectQuery) GetProjectReportLinks(reportID int64) (*[]string, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	where, args = append(where, "project_report_id = ?"), append(args, reportID)
	var projectReports []string
	err := r.conn.Select(&projectReports, `
		SELECT link
		FROM project_report_links
		WHERE `+strings.Join(where, " AND ")+`
	`, args...)
	return &projectReports, err
}

func (r *projectQuery) GetProjectRecordList(projectID int64, filter ProjectRecordFilter) (*[]ProjectRecordResponse, error) {
	where, args := []string{"pr.status > 0"}, []interface{}{}
	where, args = append(where, "pr.project_id = ?"), append(args, projectID)
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var records []ProjectRecordResponse
	err := r.conn.Select(&records, `
		SELECT pr.id, pr.project_id, pr.user_id, pr.name, pr.record_date, pr.content, pr.plan, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_records pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY pr.record_date desc, pr.updated desc
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &records, nil
}

func (r *projectQuery) GetProjectRecordCount(projectID int64) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	where, args = append(where, "project_id = ?"), append(args, projectID)
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM project_records
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *projectQuery) GetProjectRecordByID(id int64, organizationID int64) (*ProjectRecordResponse, error) {
	var record ProjectRecordResponse
	err := r.conn.Get(&record, `SELECT pr.id, pr.project_id, pr.user_id, pr.name, pr.record_date, pr.content, pr.plan, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar FROM project_records pr LEFT JOIN users u ON pr.user_id = u.id LEFT JOIN positions p ON u.position_id = p.id WHERE pr.id = ? AND pr.organization_id = ? AND pr.status > 0 limit 1`, id, organizationID)
	return &record, err
}

func (r *projectQuery) GetProjectRecordPhotos(recordID int64) (*[]string, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	where, args = append(where, "project_record_id = ?"), append(args, recordID)
	var projectRecords []string
	err := r.conn.Select(&projectRecords, `
		SELECT link
		FROM project_record_photos
		WHERE `+strings.Join(where, " AND ")+`
	`, args...)
	return &projectRecords, err
}

func (r *projectQuery) GetProjectClientUserID(id int64) (int64, error) {
	var userID int64
	err := r.conn.Get(&userID, "SELECT c.user_id FROM projects p LEFT JOIN clients c ON p.client_id = c.id WHERE p.id = ?", id)
	return userID, err
}
func (r *projectQuery) GetProjectReportViews(reportID int64) (*[]ProjectReportViewResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	where, args = append(where, "project_report_id = ?"), append(args, reportID)
	var projectReports []ProjectReportViewResponse
	err := r.conn.Select(&projectReports, `
		SELECT id, project_id, project_report_id, viewer_id, viewer_name, created
		FROM project_report_views
		WHERE `+strings.Join(where, " AND ")+`
	`, args...)
	return &projectReports, err
}

func (r *projectQuery) GetProjectReportUnreadList(userID int64) (*[]ProjectReportResponse, error) {
	var projectReports []ProjectReportResponse
	err := r.conn.Select(&projectReports, `
		SELECT pr.id, pr.user_id, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id
		WHERE pr.project_id 
		IN (SELECT project_id FROM project_members WHERE user_id = ? AND status > 0)
		AND pr.id NOT IN (SELECT project_report_id FROM project_report_views WHERE viewer_id = ? AND status > 0)
		AND pr.status > 0
		ORDER BY pr.id DESC
	`, userID, userID)
	return &projectReports, err
}
