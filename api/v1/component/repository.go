package component

import (
	"database/sql"
	"time"
)

type componentRepository struct {
	tx *sql.Tx
}

func NewComponentRepository(transaction *sql.Tx) ComponentRepository {
	return &componentRepository{
		tx: transaction,
	}
}

type ComponentRepository interface {
	//Component Management
	CreateComponent(info ComponentNew) (int64, error)
	UpdateComponent(id int64, info ComponentNew) (int64, error)
	GetComponentByID(id int64) (*Component, error)
}

func (r *componentRepository) CreateComponent(info ComponentNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO components
		(
			event_id,
			type,
			name,
			description,
			info,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.EventID, info.Type, info.Name, info.Description, info.Info, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *componentRepository) UpdateComponent(id int64, info ComponentNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update components SET 
		event_id = ?,
		type = ?,
		name = ?,
		description = ?,
		info = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.EventID, info.Type, info.Name, info.Description, info.Info, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *componentRepository) GetComponentByID(id int64) (*Component, error) {
	var res Component
	row := r.tx.QueryRow(`SELECT id, event_id, type, name, description, info, status, created, created_by, updated, updated_by FROM components WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.EventID, &res.Type, &res.Name, &res.Description, &res.Info, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
