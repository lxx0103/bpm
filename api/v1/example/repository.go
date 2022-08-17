package example

import (
	"database/sql"
	"time"
)

type exampleRepository struct {
	tx *sql.Tx
}

func NewExampleRepository(transaction *sql.Tx) ExampleRepository {
	return &exampleRepository{
		tx: transaction,
	}
}

type ExampleRepository interface {
	//Example Management
	CreateExample(ExampleNew) (int64, error)
	UpdateExample(int64, ExampleNew) (int64, error)
	GetExampleByID(int64, int64) (*Example, error)
	CheckNameExist(string, int64, int64) (int, error)
}

func (r *exampleRepository) CreateExample(info ExampleNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO examples
		(
			organization_id,
			name,
			cover,
			style,
			type,
			room,
			notes,
			description,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Cover, info.Style, info.Type, info.Room, info.Notes, info.Description, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *exampleRepository) UpdateExample(id int64, info ExampleNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update examples SET 
		name = ?,
		cover = ?,
		organization_id = ?,
		style = ?,
		type = ?,
		room = ?,
		notes = ?,
		description = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Cover, info.OrganizationID, info.Style, info.Type, info.Room, info.Notes, info.Description, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *exampleRepository) GetExampleByID(id int64, organizationID int64) (*Example, error) {
	var res Example
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, cover, style, type, room, notes, description, status, created, created_by, updated, updated_by FROM examples WHERE id = ? AND organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, cover, style, type, room, notes, description, status, created, created_by, updated, updated_by FROM examples WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Cover, &res.Style, &res.Type, &res.Room, &res.Notes, &res.Description, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *exampleRepository) CheckNameExist(name string, organizationID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM examples WHERE name = ? AND organization_id = ? AND id != ?  LIMIT 1`, name, organizationID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
