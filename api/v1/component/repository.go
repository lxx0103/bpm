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
	GetComponentByID(id int64) (*Component, error)
	GetComponentByEventID(eventID int64) (*[]Component, error)
	SaveComponent(int64, string, string) error
	CheckRequired(int64) (int, error)
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
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.EventID, info.Sort, info.Type, info.Name, info.DefaultValue, info.Required, info.Patterns, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *componentRepository) GetComponentByID(id int64) (*Component, error) {
	var res Component
	row := r.tx.QueryRow(`SELECT id, event_id, sort, component_type, name, value, default_value, required, patterns, json_data, status, created, created_by, updated, updated_by FROM event_components WHERE status > 0 AND id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.EventID, &res.Sort, &res.ComponentType, &res.Name, &res.Value, &res.DefaultValue, &res.Required, &res.Patterns, &res.JsonData, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	return &res, err
}

func (r *componentRepository) GetComponentByEventID(eventID int64) (*[]Component, error) {
	var res []Component
	rows, err := r.tx.Query(`SELECT id, required, patterns, status FROM event_components WHERE event_id = ? AND status > 0`, eventID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes Component
		err = rows.Scan(&rowRes.ID, &rowRes.Required, &rowRes.Patterns, &rowRes.Status)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *componentRepository) SaveComponent(componentID int64, value string, byUser string) error {
	_, err := r.tx.Exec(`
		Update event_components SET 
		value = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, value, 3, time.Now(), byUser, componentID)
	return err
}

func (r *componentRepository) CheckRequired(eventID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM event_components WHERE event_id = ? AND required = 1 AND status = 1`, eventID)
	err := row.Scan(&res)
	return res, err
}
