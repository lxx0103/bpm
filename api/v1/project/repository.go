package project

import (
	"database/sql"
	"time"
)

type projectRepository struct {
	tx *sql.Tx
}

func NewProjectRepository(transaction *sql.Tx) ProjectRepository {
	return &projectRepository{
		tx: transaction,
	}
}

type ProjectRepository interface {
	//Project Management
	CreateProject(info ProjectNew) (int64, error)
	UpdateProject(id int64, info ProjectNew) (int64, error)
	GetProjectByID(id int64) (*Project, error)
}

func (r *projectRepository) CreateProject(info ProjectNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO projects
		(
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *projectRepository) UpdateProject(id int64, info ProjectNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update projects SET 
		organization_id = ?,
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.OrganizationID, info.Name, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *projectRepository) GetProjectByID(id int64) (*Project, error) {
	var res Project
	row := r.tx.QueryRow(`SELECT id, organization_id, name, status, created, created_by, updated, updated_by FROM projects WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
