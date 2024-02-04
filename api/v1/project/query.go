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
		err = r.conn.Get(&project, `SELECT id, organization_id, template_id, client_id, name, type, location, longitude, latitude, checkin_distance, priority, progress, team_id, area, record_alert_day, IFNULL(last_record_date, "") as last_record_date, status, created, created_by, updated, updated_by FROM projects WHERE id = ? AND organization_id = ? AND status > 0`, id, organizationID)
	} else {
		err = r.conn.Get(&project, `SELECT id, organization_id, template_id, client_id, name, type, location, longitude, latitude, checkin_distance, priority, progress, team_id, area, record_alert_day, IFNULL(last_record_date, "") as last_record_date, status, created, created_by, updated, updated_by FROM projects WHERE id = ? AND status > 0`, id)
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
		SELECT p.id as id, p.organization_id as organization_id, o.name as organization_name, p.client_id as client_id, IFNULL(c.name, "内部流程") as client_name, p.name as name, p.type as type, p.location as location, p.longitude as longitude, p.latitude as latitude, p.checkin_distance as checkin_distance, p.priority, p.team_id as team_id, IFNULL(t.name, "") as team_name, p.area, p.record_alert_day, IFNULL(p.last_record_date, "") as last_record_date, p.status as status
		FROM projects p
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		LEFT JOIN clients c
		ON p.client_id = c.id
		LEFT JOIN teams t
		ON p.team_id = t.id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY p.id desc
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}

func (r *projectQuery) GetProjectListByCreate(userName string, organization_id int64, filter MyProjectFilter) (*[]ProjectResponse, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if filter.Status == "all" {
		where, args = append(where, "p.status > ?"), append(args, 0)
	} else {
		where, args = append(where, "p.status = ?"), append(args, 1)
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "p.type = ?"), append(args, v)
	}
	where, args = append(where, "p.created_by = ?"), append(args, userName)
	where, args = append(where, "p.organization_id = ?"), append(args, organization_id)
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var projects []ProjectResponse
	err := r.conn.Select(&projects, `
		SELECT p.id, p.organization_id, o.name as organization_name, p.template_id, IFNULL(t.name, "") as template_name, p.client_id, IFNULL(c.name, "") as client_name, p.name, p.type, p.location, p.longitude, p.latitude, p.checkin_distance, p.priority, p.team_id as team_id, IFNULL(t2.name, "") as team_name, p.area, p.record_alert_day, IFNULL(p.last_record_date, "") as last_record_date, p.progress, p.status, p.created, p.created_by, p.updated, p.updated_by
		FROM projects p
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		LEFT JOIN templates t
		ON p.template_id = t.id
		LEFT JOIN clients c
		ON p.client_id = c.id 
		LEFT JOIN teams t2
		ON p.team_id = t2.id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY p.id desc
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

func (r *projectQuery) GetProjectListByAssigned(filter AssignedProjectFilter, userID int64, positionID int64, organizationID int64) (*[]ProjectResponse, error) {
	where, args := []string{"1=1"}, []interface{}{}
	args = append(args, userID)
	if v := filter.Type; v != 0 {
		where, args = append(where, "p.type = ?"), append(args, v)
	}
	if v := filter.Status; v != 0 {
		where, args = append(where, "p.status = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "p.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.RecordStatus; v == "over" {
		where = append(where, "DATEDIFF(NOW(), IFNULL(p.last_record_date, p.created)) > p.record_alert_day")
		where = append(where, "record_alert_day > 0")
	}
	args = append(args, organizationID)
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var projects []ProjectResponse
	err := r.conn.Select(&projects, `
		SELECT p.id, p.organization_id, o.name as organization_name, p.template_id, IFNULL(t.name, "") as template_name, p.client_id, IFNULL(c.name, "") as client_name, p.name, p.type, p.location, p.longitude, p.latitude, p.checkin_distance, p.priority, p.team_id as team_id, IFNULL(t2.name, "") as team_name, p.area, p.record_alert_day, IFNULL(p.last_record_date, "") as last_record_date, p.progress, p.status, p.created, p.created_by, p.updated, p.updated_by
		FROM projects p
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		LEFT JOIN templates t
		ON p.template_id = t.id
		LEFT JOIN clients c
		ON p.client_id = c.id
		LEFT JOIN teams t2
		ON p.team_id = t2.id
		WHERE p.id IN 
		(
			SELECT project_id FROM project_members WHERE user_id = ? AND status > 0
		)
		AND `+strings.Join(where, " AND ")+`
		AND p.status > 0 AND p.organization_id = ? 
		ORDER BY p.ID DESC
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
	if v := filter.RecordStatus; v == "over" {
		where = append(where, "DATEDIFF(NOW(), IFNULL(last_record_date, created)) > record_alert_day")
		where = append(where, "record_alert_day > 0")
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

func (r *projectQuery) GetProjectListByClientID(userID int64, organization_id int64, filter MyProjectFilter) (*[]ProjectResponse, error) {
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
	var projects []ProjectResponse
	err := r.conn.Select(&projects, `
		SELECT p.id, p.organization_id, o.name as organization_name, p.template_id, IFNULL(t.name, "") as template_name, p.client_id, IFNULL(c.name, "") as client_name, p.name, p.type, p.location, p.longitude, p.latitude, p.checkin_distance, p.priority, p.team_id as team_id, IFNULL(t2.name, "") as team_name, p.area, p.record_alert_day, IFNULL(p.last_record_date, "") as last_record_date, p.progress, p.status, p.created, p.created_by, p.updated, p.updated_by
		FROM projects p
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		LEFT JOIN templates t
		ON p.template_id = t.id
		LEFT JOIN clients c
		ON p.client_id = c.id 
		LEFT JOIN teams t2
		ON p.team_id = t2.id
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
		SELECT pr.id, pr.user_id, ps.name as project_name, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id
		LEFT JOIN projects ps
		ON pr.project_id = ps.id
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id DESC
	`, args...)
	return &projectReports, err
}

func (r *projectQuery) GetProjectReportByID(id int64, organizationID int64) (*ProjectReportResponse, error) {
	var report ProjectReportResponse
	if organizationID == 0 {
		err := r.conn.Get(&report, `
		SELECT pr.id, pr.project_id, ps.name as project_name, pr.user_id, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id 
		LEFT JOIN projects ps
		ON pr.project_id = ps.id
		WHERE pr.id = ? AND pr.status > 0 limit 1
		`, id)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.conn.Get(&report, `
		SELECT pr.id, pr.project_id, ps.name as project_name, pr.user_id, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id 
		LEFT JOIN projects ps
		ON pr.project_id = ps.id
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
		SELECT pr.id, pr.project_id, ps.name, pr.user_id, pr.name, pr.record_date, pr.content, pr.plan, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_records pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id
		LEFT JOIN projects ps
		ON pr.project_id = ps.id
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
		SELECT pr.id, pr.project_id, ps.name as project_name, pr.user_id, pr.name, pr.report_date, pr.content, pr.status, pr.updated, IFNULL(u.name, "") as user_name, IFNULL(p.name, "") as position_name, u.avatar
		FROM project_reports pr
		LEFT JOIN users u
		ON pr.user_id = u.id
		LEFT JOIN positions p
		ON u.position_id = p.id
		LEFT JOIN projects ps
		ON pr.project_id = ps.id
		WHERE pr.project_id 
		IN (SELECT project_id FROM project_members WHERE user_id = ? AND status > 0)
		AND pr.id NOT IN (SELECT project_report_id FROM project_report_views WHERE viewer_id = ? AND status > 0)
		AND pr.status > 0
		ORDER BY pr.id DESC
	`, userID, userID)
	return &projectReports, err
}

func (r *projectQuery) GetActiveEvents(projectID int64) (*[]ActiveEventResponse, error) {
	var events []ActiveEventResponse
	err := r.conn.Select(&events, `
		SELECT id as event_id, name as event_name, status ,assign_type, audit_type
		FROM events 
		WHERE project_id = ? 
		AND status > 0
		AND is_active = 1
		ORDER BY sort DESC
	`, projectID)
	return &events, err
}

func (r *projectQuery) GetEventAssignPosition(eventID int64) (*[]AssignToResponse, error) {
	var events []AssignToResponse
	err := r.conn.Select(&events, `
		SELECT ea.assign_to as id , IFNULL(p.name, "") as name 
		FROM event_assigns ea 
		LEFT JOIN positions p
		ON ea.assign_to = p.id
		WHERE ea.event_id = ? 
		AND ea.status > 0
	`, eventID)
	return &events, err
}

func (r *projectQuery) GetEventAssignUser(eventID int64) (*[]AssignToResponse, error) {
	var events []AssignToResponse
	err := r.conn.Select(&events, `
		SELECT ea.assign_to as id , IFNULL(p.name, "") as name 
		FROM event_assigns ea 
		LEFT JOIN users p
		ON ea.assign_to = p.id
		WHERE ea.event_id = ? 
		AND ea.status > 0
	`, eventID)
	return &events, err
}

func (r *projectQuery) GetEventAuditPosition(eventID int64) (*[]AssignToResponse, error) {
	var auditLevel int
	err := r.conn.Get(&auditLevel, `
		SELECT audit_level
		FROM events 
		WHERE id = ? 
		AND status > 0
		LIMIT 1
	`, eventID)
	if err != nil {
		return nil, err
	}
	var auditType int
	err = r.conn.Get(&auditType, `
		SELECT audit_type
		FROM event_audits 
		WHERE event_id = ? 
		AND audit_level = ?
		AND status > 0
		LIMIT 1
	`, eventID, auditLevel)
	if err != nil {
		return nil, err
	}
	var events []AssignToResponse
	if auditType == 1 {
		err = r.conn.Select(&events, `
		SELECT ea.audit_to as id , IFNULL(p.name, "") as name 
		FROM event_audits ea 
		LEFT JOIN positions p
		ON ea.audit_to = p.id
		WHERE ea.event_id = ? 
		AND ea.audit_level = ?
		AND ea.status > 0
	`, eventID, auditLevel)
	} else if auditType == 2 {
		err = r.conn.Select(&events, `
			SELECT ea.audit_to as id , IFNULL(p.name, "") as name 
			FROM event_audits ea 
			LEFT JOIN users p
			ON ea.audit_to = p.id
			WHERE ea.event_id = ? 
			AND ea.audit_level = ?
			AND ea.status > 0
		`, eventID, auditLevel)
	}
	return &events, err
}

func (r *projectQuery) GetProjectSumByStatus(organizationID int64) (*[]ProjectSumByStatus, error) {
	var records []ProjectSumByStatus
	err := r.conn.Select(&records, `
		SELECT count(id) as sum,
		CASE 
		    WHEN status = 2 THEN '已完成'
			ELSE '进行中'
		END as status
		FROM projects
		WHERE status in (1, 2) and organization_id = ?
		GROUP BY status
	`, organizationID)
	if err != nil {
		return nil, err
	}
	return &records, nil
}

func (r *projectQuery) GetProjectSumByTeam(organizationID int64) (*[]ProjectSumByTeam, error) {
	var records []ProjectSumByTeam
	err := r.conn.Select(&records, `
		SELECT count(CASE WHEN p.status = 1 Then 1 END) as in_progress,
		count(CASE WHEN p.status = 2 Then 1 END) as completed,
		count(1) as total,
		IFNULL(t.name,"未分组") as team_name
		FROM projects p
		LEFT JOIN teams t 
		ON p.team_id = t.id
		WHERE p.status in (1, 2) and p.organization_id = ?
		GROUP BY team_id
	`, organizationID)
	if err != nil {
		return nil, err
	}
	return &records, nil
}

func (r *projectQuery) GetProjectSumByUser(organizationID int64) (*[]ProjectSumByUser, error) {
	var records []ProjectSumByUser
	err := r.conn.Select(&records, `
		SELECT count(CASE WHEN status = 1 Then 1 END) as in_progress,
		count(CASE WHEN status = 2 Then 1 END) as completed,
		count(1) as total,
		created_by as user_name
		FROM projects p
		WHERE status in (1, 2) and organization_id = ?
		GROUP BY created_by
	`, organizationID)
	if err != nil {
		return nil, err
	}
	return &records, nil
}

func (r *projectQuery) GetProjectSumByArea(organizationID int64) (*[]ProjectSumByArea, error) {
	var records []ProjectSumByArea
	err := r.conn.Select(&records, `
		SELECT count(CASE WHEN status = 1 Then 1 END) as in_progress,
		count(CASE WHEN status = 2 Then 1 END) as completed,
		count(1) as total,
		CASE when area = "" Then "未设置" ELSE area END as area_name
		FROM projects p
		WHERE status in (1, 2) and organization_id = ?
		GROUP BY area
	`, organizationID)
	if err != nil {
		return nil, err
	}
	return &records, nil
}
