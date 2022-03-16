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
	UpdateComponent(int64, Component, string) error
	GetComponentByID(id int64) (*Component, error)
	DeleteComponent(int64, string) error
}

func (r *componentRepository) CreateComponent(info ComponentNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO event_components
		(
			event_id,
			sort,
			component_type,
			name,
			default_value,
			required,
			patterns,
			json_data,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.EventID, info.Sort, info.Type, info.Name, info.DefaultValue, info.Required, info.Patterns, info.JsonData, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *componentRepository) UpdateComponent(id int64, info Component, byUser string) error {
	_, err := r.tx.Exec(`
		Update event_components SET 
		component_type = ?,
		sort = ?,
		name = ?,
		default_value = ?,
		patterns = ?,
		required = ?,
		json_data = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.ComponentType, info.Sort, info.Name, info.DefaultValue, info.Patterns, info.Required, info.JsonData, time.Now(), byUser, id)
	return err
}

func (r *componentRepository) GetComponentByID(id int64) (*Component, error) {
	var res Component
	row := r.tx.QueryRow(`SELECT id, event_id, sort, component_type, name, value, default_value, required, patterns, json_data, status, created, created_by, updated, updated_by FROM event_components WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.EventID, &res.Sort, &res.ComponentType, &res.Name, &res.Value, &res.DefaultValue, &res.Required, &res.Patterns, &res.JsonData, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	return &res, err
}

func (r *componentRepository) DeleteComponent(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update event_components SET 
		status = 2,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}
