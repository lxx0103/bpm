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
	GetEventCount(filter EventFilter) (int, error)
	GetEventList(filter EventFilter) (*[]Event, error)
}

func (r *eventQuery) GetEventByID(id int64) (*Event, error) {
	var event Event
	err := r.conn.Get(&event, "SELECT * FROM events WHERE id = ? ", id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventQuery) GetEventCount(filter EventFilter) (int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM events 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *eventQuery) GetEventList(filter EventFilter) (*[]Event, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "project_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var events []Event
	err := r.conn.Select(&events, `
		SELECT * 
		FROM events 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &events, nil
}
