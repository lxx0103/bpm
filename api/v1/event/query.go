package event

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type eventQuery struct {
	conn *sqlx.DB
}

func NewEventQuery(connection *sqlx.DB) *eventQuery {
	return &eventQuery{
		conn: connection,
	}
}

func (r *eventQuery) GetEventByID(id, organizationID int64) (*Event, error) {
	var event Event
	sql := `
	SELECT
    e.id,
    e.project_id,
    e.name,
    e.assignable,
    e.assign_type,
    e.node_id,
    e.need_audit,
    e.audit_type,
    e.complete_time,
    e.complete_user,
    e.audit_time,
    e.audit_content,
    e.audit_user,
    e.need_checkin,
    e.sort,
    e.can_review,
    IFNULL(e.deadline,"") as deadline,
    e.status,
    e.created,
    e.created_by,
    e.updated,
    e.updated_by
	FROM events e 
	LEFT JOIN projects p 
	ON e.project_id = p.id 
	WHERE e.id = ? AND e.status > 0 `
	if organizationID != 0 {
		sql += " AND p.organization_id = ? "
		err := r.conn.Get(&event, sql, id, organizationID)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.conn.Get(&event, sql, id)
		if err != nil {
			return nil, err
		}
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
		SELECT e.id, e.project_id, e.name, e.assign_type, e.node_id, e.assignable, e.need_audit, e.audit_type, e.complete_user, e.complete_time, e.audit_user, e.audit_content, e.audit_time, e.status, e.created, e.created_by, e.updated, e.updated_by, e.need_checkin, e.sort, e.can_review, IFNULL(e.deadline,"") as deadline
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
	err := r.conn.Select(&assigns, `
	SELECT ea.event_id 
	FROM event_assigns ea
	LEFT JOIN events e
	ON ea.event_id = e.id
	LEFT JOIN projects p
	ON e.project_id = p.id
	WHERE ((ea.assign_type = 2 AND ea.assign_to  = ?) OR (ea.assign_type = 1 AND ea.assign_to = ?)) AND ea.status = 1
	AND e.project_id in (SELECT project_id from project_members WHERE user_id = ? AND status > 0)
	ORDER BY p.priority asc
	`, userID, positionID, userID)
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
	sql := "SELECT e.id, e.project_id, p.name as project_name, e.name, e.complete_user, e.complete_time, e.audit_user, e.audit_time, e.audit_content, e.need_checkin, e.sort, e.status, p.priority, IFNULL(e.deadline, '') as deadline, e.can_review FROM events e LEFT JOIN projects p ON p.id = e.project_id WHERE e.id = ?"
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
	sql := "SELECT e.id, e.project_id, p.name as project_name, e.name, e.complete_user, e.complete_time, e.audit_user, e.audit_time, e.audit_content, e.need_checkin, p.priority, IFNULL(e.deadline, '') as deadline, e.sort, e.status, p.priority, e.can_review FROM events e LEFT JOIN projects p ON p.id = e.project_id WHERE e.project_id = ?  "
	if filter.Status == "all" {
		sql = sql + " AND e.status > 0"
	} else {
		sql = sql + " AND e.status in (1,3)"
	}
	sql = sql + " order by p.priority asc"
	err := r.conn.Select(&event, sql, filter.ProjectID)
	return &event, err

}

func (r *eventQuery) GetAssignedAudit(userID int64, positionID int64) ([]int64, error) {
	var assigns []int64
	err := r.conn.Select(&assigns, `
		SELECT event_id FROM event_audits ea 
		LEFT JOIN events e
		ON ea.event_id = e.id
		WHERE ((ea.audit_type = 2 AND ea.audit_to  = ?) OR (ea.audit_type = 1 AND ea.audit_to = ?)) 
		AND e.project_id IN (SELECT project_id from project_members WHERE user_id = ? AND status > 0)
		AND ea.status = 1
	`, userID, positionID, userID)
	return assigns, err
}

func (r *eventQuery) GetAssignedAuditByID(id int64, status string) (*MyEvent, error) {
	var event MyEvent
	sql := "SELECT e.id, e.project_id, p.name as project_name, e.name, e.complete_user, e.complete_time, e.audit_user, e.audit_time, e.audit_content, e.need_checkin, e.sort, e.status, p.priority, IFNULL(e.deadline, '') as deadline, e.can_review FROM events e LEFT JOIN projects p ON p.id = e.project_id WHERE e.id = ?"
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

func (r *eventQuery) GetCheckinCount(filter CheckinFilter) (int, error) {
	where, args := []string{"ec.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "ec.user_name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "ec.user_id = ?"), append(args, v)
	}
	if v := filter.EventID; v != 0 {
		where, args = append(where, "ec.event_id = ?"), append(args, v)
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "e.project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	if v := filter.From; v != "" {
		where, args = append(where, "ec.checkin_time >= ?"), append(args, v+" 00:00:00")
	}
	if v := filter.To; v != "" {
		where, args = append(where, "ec.checkin_time <= ?"), append(args, v+" 23:59:59")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM event_checkins ec
		LEFT JOIN events e
		ON ec.event_id = e.id
		LEFT JOIN projects p 
		ON e.project_id = p.id
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *eventQuery) GetCheckinList(filter CheckinFilter) (*[]CheckinResponse, error) {
	where, args := []string{"ec.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "ec.user_name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "ec.user_id = ?"), append(args, v)
	}
	if v := filter.EventID; v != 0 {
		where, args = append(where, "ec.event_id = ?"), append(args, v)
	}
	if v := filter.ProjectID; v != 0 {
		where, args = append(where, "e.project_id = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	if v := filter.From; v != "" {
		where, args = append(where, "ec.checkin_time >= ?"), append(args, v+" 00:00:00")
	}
	if v := filter.To; v != "" {
		where, args = append(where, "ec.checkin_time <= ?"), append(args, v+" 23:59:59")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var checkins []CheckinResponse
	err := r.conn.Select(&checkins, `
		SELECT ec.user_name as name, e.project_id as project_id, p.name as project_name, ec.event_id as event_id, e.name as event_name, p.organization_id as organization_id, o.name as organization_name, ec.checkin_type as checkin_type, ec.checkin_time as checkin_time, ec.longitude as longitude, ec.latitude as latitude, ec.distance as distance
		FROM event_checkins ec
		LEFT JOIN events e
		ON ec.event_id = e.id
		LEFT JOIN projects p 
		ON e.project_id = p.id
		LEFT JOIN organizations o
		ON p.organization_id = o.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &checkins, nil
}

func (r *eventQuery) GetAuditHistoryList(eventID int64) (*[]EventAuditHistoryResponse, error) {
	var historys []EventAuditHistoryResponse
	err := r.conn.Select(&historys, `
		SELECT id, event_id, audit_user, audit_time, audit_content, status
		FROM event_audit_historys
		WHERE event_id = ? AND status > 0
		ORDER BY audit_time asc
	`, eventID)
	return &historys, err
}

func (r *eventQuery) GetReviewList(eventID int64) (*[]EventReviewResponse, error) {
	var reviews []EventReviewResponse
	err := r.conn.Select(&reviews, `
		SELECT id, event_id, result, content, link, status, created, handle_user, handle_content, handle_time
		FROM event_reviews
		WHERE event_id = ? AND status > 0
		ORDER BY id desc
	`, eventID)
	return &reviews, err
}
