package event

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type eventQuery struct {
	conn *sqlx.DB
}

func NewEventQuery(connection *sqlx.DB) EventQuery {
	return &eventQuery{
		conn: connection,
	}
}

type EventQuery interface {
	//Event Management
	GetEventByID(id int64) (*Event, error)
	GetAssignsByEventID(int64) (*[]EventAssign, error)
	GetPresByEventID(int64) (*[]EventPre, error)
	GetAuditsByEventID(int64) (*[]EventAudit, error)
	GetEventCount(EventFilter, int64) (int, error)
	GetEventList(EventFilter, int64) (*[]Event, error)
	//WX
	GetAssigned(int64, int64) ([]int64, error)
	GetAssignedAudit(int64, int64) ([]int64, error)
	CheckActive(int64) (bool, error)
	GetAssignedEventByID(int64, string) (*MyEvent, error)
	GetProjectEvent(MyEventFilter) (*[]MyEvent, error)
	GetAssignedAuditByID(int64, string) (*MyEvent, error)
}

func (r *eventQuery) GetEventByID(id int64) (*Event, error) {
	var event Event
	err := r.conn.Get(&event, "SELECT * FROM events WHERE id = ? AND status > 0 ", id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventQuery) GetEventCount(filter EventFilter, organizationID int64) (int, error) {
	where, args := []string{"e.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "e.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "e.project_id = ?"), append(args, v)
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM events e
		LEFT JOIN projects p
		ON e.project_id = p.id
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *eventQuery) GetEventList(filter EventFilter, organizationID int64) (*[]Event, error) {
	where, args := []string{"e.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "e.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "e.project_id = ?"), append(args, v)
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var events []Event
	err := r.conn.Select(&events, `
		SELECT e.* 
		FROM events e
		LEFT JOIN projects p
		ON e.project_id = p.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &events, nil
}

func (r *eventQuery) GetAssignsByEventID(eventID int64) (*[]EventAssign, error) {
	var assigns []EventAssign
	err := r.conn.Select(&assigns, "SELECT * FROM event_assigns WHERE event_id = ? AND status = ?", eventID, 1)
	if err != nil {
		return nil, err
	}
	return &assigns, nil
}

func (r *eventQuery) GetPresByEventID(eventID int64) (*[]EventPre, error) {
	var pres []EventPre
	err := r.conn.Select(&pres, "SELECT * FROM event_pres WHERE event_id = ? AND status = ?", eventID, 1)
	if err != nil {
		return nil, err
	}
	return &pres, nil
}

func (r *eventQuery) GetAuditsByEventID(eventID int64) (*[]EventAudit, error) {
	var pres []EventAudit
	err := r.conn.Select(&pres, "SELECT * FROM event_audits WHERE event_id = ? AND status = ?", eventID, 1)
	if err != nil {
		return nil, err
	}
	return &pres, nil
}

func (r *eventQuery) GetAssigned(userID int64, positionID int64) ([]int64, error) {
	var assigns []int64
	err := r.conn.Select(&assigns, "SELECT event_id FROM event_assigns WHERE ((assign_type = 2 AND assign_to  = ?) OR (assign_type = 1 AND assign_to = ?)) AND status = 1", userID, positionID)
	return assigns, err
}

func (r *eventQuery) CheckActive(eventID int64) (bool, error) {
	var activePreCount int
	err := r.conn.Get(&activePreCount, `
		SELECT count(1) from event_pres ep
		LEFT JOIN events e
		ON ep.pre_id = e.id 
		WHERE ep.status > 0  
		AND ep.event_id = ?
		AND e.status not in (-1, 9)`, eventID)
	if err != nil {
		return false, err
	}
	if activePreCount == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *eventQuery) GetAssignedEventByID(id int64, status string) (*MyEvent, error) {
	var event MyEvent
	sql := "SELECT e.id, e.project_id, p.name as project_name, e.name, e.complete_user, e.complete_time, e.audit_user, e.audit_time, e.audit_content, e.need_checkin, e.status FROM events e LEFT JOIN projects p ON p.id = e.project_id WHERE e.id = ?"
	if status == "all" {
		sql = sql + " AND e.status > 0"
	} else {
		sql = sql + " AND e.status in (1,3)"
	}
	err := r.conn.Get(&event, sql, id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventQuery) GetProjectEvent(filter MyEventFilter) (*[]MyEvent, error) {
	var event []MyEvent
	sql := "SELECT e.id, e.project_id, p.name as project_name, e.name, e.complete_user, e.complete_time, e.audit_user, e.audit_time, e.audit_content, e.need_checkin, e.status FROM events e LEFT JOIN projects p ON p.id = e.project_id WHERE e.project_id = ?  "
	if filter.Status == "all" {
		sql = sql + " AND e.status > 0"
	} else {
		sql = sql + " AND e.status in (1,3)"
	}
	err := r.conn.Select(&event, sql, filter.ProjectID)
	return &event, err

}

func (r *eventQuery) GetAssignedAudit(userID int64, positionID int64) ([]int64, error) {
	var assigns []int64
	err := r.conn.Select(&assigns, "SELECT event_id FROM event_audits WHERE ((audit_type = 2 AND audit_to  = ?) OR (audit_type = 1 AND audit_to = ?)) AND status = 1", userID, positionID)
	return assigns, err
}

func (r *eventQuery) GetAssignedAuditByID(id int64, status string) (*MyEvent, error) {
	var event MyEvent
	sql := "SELECT e.id, e.project_id, p.name as project_name, e.name, e.complete_user, e.complete_time, e.audit_user, e.audit_time, e.audit_content, e.need_checkin, e.status FROM events e LEFT JOIN projects p ON p.id = e.project_id WHERE e.id = ?"
	if status == "all" {
		sql = sql + " AND e.status > 0"
	} else {
		sql = sql + " AND e.status = 2"
	}
	err := r.conn.Get(&event, sql, id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
