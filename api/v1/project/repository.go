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
