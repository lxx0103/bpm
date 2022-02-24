package event

import (
	"database/sql"
	"errors"
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
	CreateEventAssign(int64, int, []int64, string) error
	DeleteEventAssign(int64, string) error
	UpdateEvent(int64, Event, string) (int64, error)
	GetEventByID(int64, int64) (*Event, error)
	GetAssignsByEventID(int64) (*[]EventAssign, error)
	CheckProjectExist(int64, int64) (int, error)
	CheckNameExist(string, int64, int64) (int, error)
}

func (r *eventRepository) CreateEvent(info EventNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO events
		(
			project_id,
			name,
			pre_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.ProjectID, info.Name, info.PreID, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *eventRepository) CreateEventAssign(eventID int64, assignType int, assignTo []int64, user string) error {
	for i := 0; i < len(assignTo); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM event_assigns WHERE event_id = ? AND assign_type = ? AND assign_to = ? AND status = 1  LIMIT 1`, eventID, assignType, assignTo[i])
		err := row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist != 0 {
			msg := "指派对象有重复"
			return errors.New(msg)
		}
		_, err = r.tx.Exec(`
			INSERT INTO event_assigns
			(
				event_id,
				assign_type,
				assign_to,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, eventID, assignType, assignTo[i], 1, time.Now(), user, time.Now(), user)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *eventRepository) UpdateEvent(id int64, info Event, byUser string) (int64, error) {
	result, err := r.tx.Exec(`
		Update events SET 
		name = ?,
		pre_id = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.PreID, info.Status, time.Now(), byUser, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *eventRepository) DeleteEventAssign(event_id int64, user string) error {
	_, err := r.tx.Exec(`
		Update event_assigns SET
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE event_id = ?
	`, 2, time.Now(), user, event_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetEventByID(id int64, organizationID int64) (*Event, error) {
	var res Event
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT e.id, e.project_id, e.name, e.pre_id, e.status, e.created, e.created_by, e.updated, e.updated_by FROM events e LEFT JOIN projects p ON e.project_id = p.id  WHERE e.id = ? AND p.organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, project_id, name, pre_id, status, created, created_by, updated, updated_by FROM events WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.ProjectID, &res.Name, &res.PreID, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *eventRepository) CheckProjectExist(projectID int64, organizationID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM projects WHERE id = ? AND organization_id = ?  LIMIT 1`, projectID, organizationID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *eventRepository) CheckNameExist(name string, projectID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM events WHERE name = ? AND project_id = ? AND id != ?  LIMIT 1`, name, projectID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *eventRepository) GetAssignsByEventID(eventID int64) (*[]EventAssign, error) {
	var res []EventAssign
	rows, err := r.tx.Query(`SELECT id, event_id, assign_type, assign_to, status, created, created_by, updated, updated_by FROM event_assigns WHERE event_id = ? AND status = ? `, eventID, 1)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes EventAssign
		err = rows.Scan(&rowRes.ID, &rowRes.EventID, &rowRes.AssignType, &rowRes.AssignTo, &rowRes.Status, &rowRes.Created, &rowRes.CreatedBy, &rowRes.Updated, &rowRes.UpdatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}
