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
	GetAssignsByEventID(int64) (*[]EventAssign, error)
	CreateEventAudit(int64, int, []int64, string) error
	DeleteEventAudit(int64, string) error
	GetAuditsByEventID(int64) (*[]EventAudit, error)
	CreateEventPre(int64, []int64, string) error
	DeleteEventPre(int64, string) error
	GetPresByEventID(int64) (*[]EventPre, error)
	UpdateEvent(int64, Event, string) error
	GetEventByID(int64, int64) (*Event, error)
	DeleteEventByProjectID(int64, string) error
	CheckProjectExist(int64, int64) (int, error)
	CheckNameExist(string, int64, int64) (int, error)
	GetEventsByProjectID(int64) (*[]Event, error)
	GetEventIDByProjectAndNode(int64, int64) (int64, error)
	CheckAssign(int64, int64, int64) (int, error)
	CompleteEvent(int64, string) error
}

func (r *eventRepository) CreateEvent(info EventNew) (int64, error) {
	if info.AssignType == 3 {
		info.AssignType = 2
	}
	result, err := r.tx.Exec(`
		INSERT INTO events
		(
			project_id,
			node_id,
			name,
			assign_type,
			assignable,
			need_audit,
			audit_type,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.ProjectID, info.NodeID, info.Name, info.AssignType, info.Assignable, info.NeedAudit, info.AuditType, 1, time.Now(), info.User, time.Now(), info.User)
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
		if assignType == 3 {
			assignType = 2
		}
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM event_assigns WHERE event_id = ? AND assign_type = ? AND assign_to = ? AND status > 0  LIMIT 1`, eventID, assignType, assignTo[i])
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

func (r *eventRepository) DeleteEventAssign(event_id int64, user string) error {
	_, err := r.tx.Exec(`
		Update event_assigns SET
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE event_id = ?
	`, time.Now(), user, event_id)
	if err != nil {
		return err
	}
	return nil
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

func (r *eventRepository) UpdateEvent(id int64, info Event, byUser string) error {
	_, err := r.tx.Exec(`
		Update events SET 
		assign_type = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.AssignType, time.Now(), byUser, id)
	return err
}

func (r *eventRepository) GetEventByID(id int64, organizationID int64) (*Event, error) {
	var res Event
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT e.id, e.project_id, e.name, e.assignable, e.assign_type, e.status, e.created, e.created_by, e.updated, e.updated_by FROM events e LEFT JOIN projects p ON e.project_id = p.id  WHERE e.id = ? AND p.organization_id = ? AND e.status > 0 LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, project_id, name, assignable, assign_type, status, created, created_by, updated, updated_by FROM events WHERE id = ? AND status > 0 LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.ProjectID, &res.Name, &res.Assignable, &res.AssignType, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
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
	row := r.tx.QueryRow(`SELECT count(1) FROM events WHERE name = ? AND project_id = ? AND id != ? AND status > 0  LIMIT 1`, name, projectID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *eventRepository) CreateEventPre(eventID int64, preIDs []int64, user string) error {
	for i := 0; i < len(preIDs); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM event_pres WHERE event_id = ? AND pre_id = ? AND status > 0  LIMIT 1`, eventID, preIDs[i])
		err := row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist != 0 {
			msg := "前置事件有重复"
			return errors.New(msg)
		}
		_, err = r.tx.Exec(`
			INSERT INTO event_pres
			(
				event_id,
				pre_id,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, eventID, preIDs[i], 1, time.Now(), user, time.Now(), user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *eventRepository) DeleteEventPre(event_id int64, user string) error {
	_, err := r.tx.Exec(`
		Update event_pres SET
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE event_id = ?
	`, time.Now(), user, event_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetPresByEventID(eventID int64) (*[]EventPre, error) {
	var res []EventPre
	rows, err := r.tx.Query(`SELECT id, event_id, pre_id, status, created, created_by, updated, updated_by FROM event_pres WHERE event_id = ? AND status > 0 `, eventID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes EventPre
		err = rows.Scan(&rowRes.ID, &rowRes.EventID, &rowRes.PreID, &rowRes.Status, &rowRes.Created, &rowRes.CreatedBy, &rowRes.Updated, &rowRes.UpdatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *eventRepository) DeleteEventByProjectID(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update events SET status = -1, updated = ?,updated_by = ? WHERE project_id = ?`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		Update event_pres ep 
		LEFT JOIN events e 
		ON ep.event_id = e.id 
		SET	ep.status = -1,
		ep.updated = ?,
		ep.updated_by = ?
		WHERE e.project_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		Update event_assigns ea 
		LEFT JOIN events e 
		ON ea.event_id = e.id 
		SET	ea.status = -1,
		ea.updated = ?,
		ea.updated_by = ?
		WHERE e.project_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		Update event_components ec 
		LEFT JOIN events e 
		ON ec.event_id = e.id 
		SET	ec.status = -1,
		ec.updated = ?,
		ec.updated_by = ?
		WHERE e.project_id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *eventRepository) GetEventsByProjectID(projectID int64) (*[]Event, error) {
	var res []Event
	rows, err := r.tx.Query(`SELECT id, project_id, node_id, name, assign_type, assignable, need_audit, audit_type FROM events WHERE project_id = ? AND status > 0`, projectID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes Event
		err = rows.Scan(&rowRes.ID, &rowRes.ProjectID, &rowRes.NodeID, &rowRes.Name, &rowRes.AssignType, &rowRes.Assignable, &rowRes.NeedAudit, &rowRes.AuditType)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *eventRepository) GetEventIDByProjectAndNode(projectID int64, nodeID int64) (int64, error) {
	var res int64
	row := r.tx.QueryRow(`SELECT id FROM events WHERE project_id = ? AND node_id = ? AND status > 0 LIMIT 1`, projectID, nodeID)
	err := row.Scan(&res)
	return res, err
}

func (r *eventRepository) CheckAssign(eventID int64, userID int64, positionID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM event_assigns WHERE event_id = ? AND ( ( assign_type = 1 AND assign_to = ? ) OR ( assign_type = 2 and assign_to = ? ) ) AND status > 0  LIMIT 1`, eventID, positionID, userID)
	err := row.Scan(&res)
	return res, err
}

func (r *eventRepository) CompleteEvent(eventID int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update events SET 
		complete_user = ?,
		complete_time = ?,
		status = 3,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, byUser, time.Now().Format("2006-01-02 15:04:05"), time.Now(), byUser, eventID)
	return err
}

func (r *eventRepository) CreateEventAudit(eventID int64, auditType int, auditTo []int64, user string) error {
	for i := 0; i < len(auditTo); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM event_audits WHERE event_id = ? AND audit_type = ? AND audit_to = ? AND status > 0  LIMIT 1`, eventID, auditType, auditTo[i])
		err := row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist != 0 {
			msg := "指派对象有重复"
			return errors.New(msg)
		}
		_, err = r.tx.Exec(`
			INSERT INTO event_audits
			(
				event_id,
				audit_type,
				audit_to,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, eventID, auditType, auditTo[i], 1, time.Now(), user, time.Now(), user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *eventRepository) DeleteEventAudit(event_id int64, user string) error {
	_, err := r.tx.Exec(`
		Update event_audits SET
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE event_id = ?
	`, time.Now(), user, event_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetAuditsByEventID(eventID int64) (*[]EventAudit, error) {
	var res []EventAudit
	rows, err := r.tx.Query(`SELECT id, event_id, audit_type, audit_to, status, created, created_by, updated, updated_by FROM event_audits WHERE event_id = ? AND status = ? `, eventID, 1)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes EventAudit
		err = rows.Scan(&rowRes.ID, &rowRes.EventID, &rowRes.AuditType, &rowRes.AuditTo, &rowRes.Status, &rowRes.Created, &rowRes.CreatedBy, &rowRes.Updated, &rowRes.UpdatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}
