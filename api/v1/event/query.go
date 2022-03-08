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
	GetEventCount(EventFilter, int64) (int, error)
	GetEventList(EventFilter, int64) (*[]Event, error)
}

func (r *eventQuery) GetEventByID(id int64) (*Event, error) {
	var event Event
	err := r.conn.Get(&event, "SELECT * FROM events WHERE id = ? AND status = 1 ", id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventQuery) GetEventCount(filter EventFilter, organizationID int64) (int, error) {
	where, args := []string{"status = 1"}, []interface{}{}
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
	where, args := []string{"status = 1"}, []interface{}{}
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
