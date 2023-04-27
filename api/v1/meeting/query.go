package meeting

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type meetingQuery struct {
	conn *sqlx.DB
}

func NewMeetingQuery(connection *sqlx.DB) *meetingQuery {
	return &meetingQuery{
		conn: connection,
	}
}

func (r *meetingQuery) GetMeetingByID(id int64, organizationID int64) (*MeetingResponse, error) {
	var meeting MeetingResponse
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&meeting, `
		SELECT m.id, m.name, m.status, m.organization_id, o.name as organization_name, m.date, m.content, m.file, m.user_id
		FROM meetings m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE m.id = ? 
		AND m.organization_id = ? 
		AND m.status > 0
		`, id, organizationID)
	} else {
		err = r.conn.Get(&meeting, `
		SELECT m.id, m.name, m.status, m.organization_id, o.name as organization_name, m.date, m.content, m.file, m.user_id
		FROM meetings m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE m.id = ? 
		AND m.status > 0
		`, id)
	}
	return &meeting, err
}

func (r *meetingQuery) GetMeetingCount(filter MeetingFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM meetings 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *meetingQuery) GetMeetingList(filter MeetingFilter) (*[]MeetingResponse, error) {
	where, args := []string{"m.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "m.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "m.organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var meetings []MeetingResponse
	err := r.conn.Select(&meetings, `
		SELECT m.id, m.name, m.status, m.organization_id, o.name as organization_name, m.date, m.content, m.file, m.user_id
		FROM meetings m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &meetings, nil
}
