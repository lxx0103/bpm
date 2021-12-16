package event

import (
	"database/sql"
	"time"
)

type eventRepository struct {
	tx *sql.Tx
}

func NewEventRepository(transaction *sql.Tx) EventRepository {
	return &eventRepository{
		tx: transaction,
	}
}

type EventRepository interface {
	//Event Management
	CreateEvent(info EventNew) (int64, error)
	UpdateEvent(id int64, info EventNew) (int64, error)
	GetEventByID(id int64) (*Event, error)
}

func (r *eventRepository) CreateEvent(info EventNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO events
		(
			project_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, info.ProjectID, info.Name, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *eventRepository) UpdateEvent(id int64, info EventNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update events SET 
		project_id = ?,
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.ProjectID, info.Name, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *eventRepository) GetEventByID(id int64) (*Event, error) {
	var res Event
	row := r.tx.QueryRow(`SELECT id, project_id, name, status, created, created_by, updated, updated_by FROM events WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.ProjectID, &res.Name, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
