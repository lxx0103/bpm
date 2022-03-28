package template

import (
	"database/sql"
	"time"
)

type templateRepository struct {
	tx *sql.Tx
}

func NewTemplateRepository(transaction *sql.Tx) TemplateRepository {
	return &templateRepository{
		tx: transaction,
	}
}

type TemplateRepository interface {
	//Template Management
	CreateTemplate(TemplateNew) (int64, error)
	UpdateTemplate(int64, Template, string) error
	GetTemplateByID(int64) (*Template, error)
	CheckNameExist(string, int64, int64) (int, error)
	DeleteTemplate(int64, string) error
}

func (r *templateRepository) CreateTemplate(info TemplateNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO templates
		(
			organization_id,
			name,
			event_json,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.EventJson, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *templateRepository) UpdateTemplate(id int64, info Template, byUser string) error {
	_, err := r.tx.Exec(`
		Update templates SET 
		name = ?,
		status = ?,
		event_json = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Status, info.EventJson, time.Now(), byUser, id)
	return err
}

func (r *templateRepository) GetTemplateByID(id int64) (*Template, error) {
	var res Template
	row := r.tx.QueryRow(`SELECT id, organization_id, name, event_json, status, created, created_by, updated, updated_by FROM templates WHERE status > 0 AND id = ? LIMIT 1`, id)

	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.EventJson, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	return &res, err
}

func (r *templateRepository) CheckNameExist(name string, organizationID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM templates WHERE status > 0 AND name = ? AND organization_id = ? AND id != ? LIMIT 1`, name, organizationID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *templateRepository) DeleteTemplate(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update templates SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}
