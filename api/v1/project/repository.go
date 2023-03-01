package project

import (
	"database/sql"
	"time"
)

type projectRepository struct {
	tx *sql.Tx
}

func NewProjectRepository(transaction *sql.Tx) *projectRepository {
	return &projectRepository{
		tx: transaction,
	}
}

func (r *projectRepository) CreateProject(info ProjectNew, organizationID int64) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO projects
		(
			organization_id,
			template_id,
			client_id,
			name,
			type,
			location,
			longitude,
			latitude,
			checkin_distance,
			priority,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, organizationID, info.TemplateID, info.ClientID, info.Name, info.Type, info.Location, info.Longitude, info.Latitude, info.CheckinDistance, info.Priority, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *projectRepository) UpdateProject(id int64, info Project, byUser string) error {
	_, err := r.tx.Exec(`
		Update projects SET 
		name = ?,
		client_id = ?,
		location = ?,
		longitude = ?,
		latitude = ?,
		checkin_distance = ?,
		priority = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.ClientID, info.Location, info.Longitude, info.Latitude, info.CheckinDistance, info.Priority, time.Now(), byUser, id)
	return err
}

func (r *projectRepository) GetProjectByID(id int64, organizationID int64) (*Project, error) {
	var res Project
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, type, location, longitude, latitude, checkin_distance, status, created, created_by, updated, updated_by FROM projects WHERE id = ? AND organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, type, location, longitude, latitude, checkin_distance, status, created, created_by, updated, updated_by FROM projects WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Type, &res.Location, &res.Longitude, &res.Latitude, &res.CheckinDistance, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *projectRepository) CheckNameExist(name string, organizationID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM projects WHERE name = ? AND organization_id = ? AND id != ? AND status > 0 LIMIT 1`, name, organizationID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *projectRepository) DeleteProject(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update projects SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *projectRepository) CreateProjectReport(info ProjectReport) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO project_reports
		(
			organization_id,
			project_id,
			client_id,
			user_id,
			report_date,
			name,
			content,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ProjectID, info.ClientID, info.UserID, info.ReportDate, info.Name, info.Content, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *projectRepository) CreateProjectReportLink(info ProjectReportLink) error {
	_, err := r.tx.Exec(`
		INSERT INTO project_report_links
		(
			organization_id,
			project_id,
			project_report_id,
			link,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ProjectID, info.ProjectReportID, info.Link, info.Status, time.Now(), info.CreatedBy, time.Now(), info.UpdatedBy)
	return err
}

func (r *projectRepository) GetProjectReportByID(id int64, organizationID int64) (*ProjectReportResponse, error) {
	var res ProjectReportResponse
	row := r.tx.QueryRow(`SELECT id, project_id, name, report_date, content, updated, status, user_id FROM project_reports WHERE id = ? AND organization_id = ? AND status > 0 LIMIT 1`, id, organizationID)

	err := row.Scan(&res.ID, &res.ProjectID, &res.Name, &res.ReportDate, &res.Content, &res.Updated, &res.Status, &res.UserID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *projectRepository) DeleteProjectReport(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update project_reports SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}
func (r *projectRepository) DeleteProjectReportLinks(reportID int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update project_report_links SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE project_report_id = ?
	`, time.Now(), byUser, reportID)
	return err
}
func (r *projectRepository) DeleteProjectReportViews(reportID int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update project_report_views SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE project_report_id = ?
	`, time.Now(), byUser, reportID)
	return err
}

func (r *projectRepository) UpdateProjectReport(id int64, info ProjectReport) error {
	_, err := r.tx.Exec(`
		Update project_reports SET
		report_date = ?,
		name = ?,
		content = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.ReportDate, info.Name, info.Content, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *projectRepository) CreateProjectRecord(info ProjectRecord) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO project_records
		(
			organization_id,
			project_id,
			client_id,
			user_id,
			record_date,
			name,
			content,
			plan,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ProjectID, info.ClientID, info.UserID, info.RecordDate, info.Name, info.Content, info.Plan, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *projectRepository) CreateProjectRecordPhoto(info ProjectRecordPhoto) error {
	_, err := r.tx.Exec(`
		INSERT INTO project_record_photos
		(
			organization_id,
			project_id,
			project_record_id,
			link,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ProjectID, info.ProjectRecordID, info.Link, info.Status, time.Now(), info.CreatedBy, time.Now(), info.UpdatedBy)
	return err
}

func (r *projectRepository) GetProjectRecordByID(id int64, organizationID int64) (*ProjectRecordResponse, error) {
	var res ProjectRecordResponse
	row := r.tx.QueryRow(`SELECT id, project_id, name, record_date, content, updated, status, user_id FROM project_records WHERE id = ? AND organization_id = ? AND status > 0 LIMIT 1`, id, organizationID)

	err := row.Scan(&res.ID, &res.ProjectID, &res.Name, &res.RecordDate, &res.Content, &res.Updated, &res.Status, &res.UserID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *projectRepository) DeleteProjectRecord(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update project_records SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}
func (r *projectRepository) DeleteProjectRecordPhotos(recordID int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update project_record_photos SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE project_record_id = ?
	`, time.Now(), byUser, recordID)
	return err
}

func (r *projectRepository) UpdateProjectRecord(id int64, info ProjectRecord) error {
	_, err := r.tx.Exec(`
		Update project_records SET
		record_date = ?,
		name = ?,
		content = ?,
		plan = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.RecordDate, info.Name, info.Content, info.Plan, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *projectRepository) CreateProjectReportView(info ProjectReportView) error {
	_, err := r.tx.Exec(`
		INSERT INTO project_report_views
		(
			organization_id,
			project_id,
			project_report_id,
			viewer_id,
			viewer_name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ProjectID, info.ProjectReportID, info.ViewerID, info.ViewerName, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *projectRepository) GetProjectReportView(id int64) (*[]ProjectReportViewResponse, error) {
	var res []ProjectReportViewResponse
	rows, err := r.tx.Query(`SELECT id, project_id, project_report_id, viewer_id, viewer_name, created FROM project_report_views WHERE project_report_id = ? AND status > 0`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes ProjectReportViewResponse
		err = rows.Scan(&rowRes.ID, &rowRes.ProjectID, &rowRes.ProjectReportID, &rowRes.ViewerID, &rowRes.ViewerName, &rowRes.Created)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}
func (r *projectRepository) UpdateProjectReportStatus(id int64, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update project_reports SET
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, status, time.Now(), byUser, id)
	return err
}
func (r *projectRepository) CheckViewExist(reportID, userID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM project_report_views WHERE project_report_id = ? AND viewer_id = ? AND status > 0 LIMIT 1`, reportID, userID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
