package event

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type eventRepository struct {
	tx *sql.Tx
}

func NewEventRepository(transaction *sql.Tx) *eventRepository {
	return &eventRepository{
		tx: transaction,
	}
}

func (r *eventRepository) CreateEvent(info EventNew) (int64, error) {
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
			audit_level,
			need_checkin,
			sort,
			can_review,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.ProjectID, info.NodeID, info.Name, info.AssignType, info.Assignable, info.NeedAudit, info.AuditType, 1, info.NeedCheckin, info.Sort, info.CanReview, 1, time.Now(), info.User, time.Now(), info.User)
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
		need_audit = ?,
		audit_type = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.AssignType, info.NeedAudit, info.AuditType, time.Now(), byUser, id)
	return err
}

func (r *eventRepository) GetEventByID(id int64, organizationID int64) (*Event, error) {
	var res Event
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT e.id, e.project_id, e.name, e.assignable, e.assign_type, e.need_audit, e.audit_level, e.audit_type, e.audit_content, e.audit_time, e.audit_user, e.need_checkin, e.sort, e.can_review, IFNULL(e.deadline,"") as deadline, e.status, e.created, e.created_by, e.updated, e.updated_by FROM events e LEFT JOIN projects p ON e.project_id = p.id  WHERE e.id = ? AND p.organization_id = ? AND e.status > 0 LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, project_id, name, assignable, assign_type, need_audit, audit_level, audit_type, audit_content, audit_time, audit_user, need_checkin, sort, can_review, IFNULL(deadline,"") as deadline, status, created, created_by, updated, updated_by FROM events WHERE id = ? AND status > 0 LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.ProjectID, &res.Name, &res.Assignable, &res.AssignType, &res.NeedAudit, &res.AuditLevel, &res.AuditType, &res.AuditContent, &res.AuditTime, &res.AuditUser, &res.NeedCheckin, &res.Sort, &res.CanReview, &res.Deadline, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		fmt.Println(err)
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
	rows, err := r.tx.Query(`SELECT id, project_id, node_id, name, assign_type, assignable, need_audit, audit_type, can_review FROM events WHERE project_id = ? AND status > 0`, projectID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes Event
		err = rows.Scan(&rowRes.ID, &rowRes.ProjectID, &rowRes.NodeID, &rowRes.Name, &rowRes.AssignType, &rowRes.Assignable, &rowRes.NeedAudit, &rowRes.AuditType, &rowRes.CanReview)
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

func (r *eventRepository) CheckAudit(eventID int64, userID int64, positionID int64, auditLevel int) (int, error) {
	var res int
	if auditLevel != 0 {
		row := r.tx.QueryRow(`SELECT count(1) FROM event_audits WHERE event_id = ? AND ( ( audit_type = 1 AND audit_to = ? ) OR ( audit_type = 2 and audit_to = ? ) ) AND status > 0 AND audit_level = ? LIMIT 1`, eventID, positionID, userID, auditLevel)
		err := row.Scan(&res)
		if err != nil {
			return 0, err
		}
	} else {
		row := r.tx.QueryRow(`SELECT count(1) FROM event_audits WHERE event_id = ? AND ( ( audit_type = 1 AND audit_to = ? ) OR ( audit_type = 2 and audit_to = ? ) ) AND status > 0  LIMIT 1`, eventID, positionID, userID)
		err := row.Scan(&res)
		if err != nil {
			return 0, err
		}
	}
	return res, nil
}
func (r *eventRepository) CompleteEvent(eventID int64, byUser string) (int64, error) {
	_, err := r.tx.Exec(`
		Update events SET 
		complete_user = ?,
		complete_time = ?,
		status = 2,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, byUser, time.Now().Format("2006-01-02 15:04:05"), time.Now(), byUser, eventID)
	if err != nil {
		return 0, err
	}
	res, err := r.tx.Exec(`
		INSERT INTO event_historys
		(
			event_id,
			history_type,
			audit_user,
			audit_content,
			audit_time,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, eventID, "完成事件", byUser, "", time.Now().Format("2006-01-02 15:04:05"), 1, time.Now(), byUser, time.Now(), byUser)
	if err != nil {
		return 0, err
	}
	historyID, err := res.LastInsertId()
	return historyID, err
}

func (r *eventRepository) CreateEventAudit(eventID int64, auditType int, auditInfo NodeAudit, user string) error {
	for i := 0; i < len(auditInfo.AuditTo); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM event_audits WHERE event_id = ? AND audit_level = ? AND audit_type = ? AND audit_to = ? AND status > 0  LIMIT 1`, eventID, auditInfo.AuditLevel, auditType, auditInfo.AuditTo[i])
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
				audit_level,
				audit_type,
				audit_to,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, eventID, auditInfo.AuditLevel, auditType, auditInfo.AuditTo[i], 1, time.Now(), user, time.Now(), user)
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
	rows, err := r.tx.Query(`SELECT id, event_id, audit_level, audit_type, audit_to, status, created, created_by, updated, updated_by FROM event_audits WHERE event_id = ? AND status = ? `, eventID, 1)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes EventAudit
		err = rows.Scan(&rowRes.ID, &rowRes.EventID, &rowRes.AuditLevel, &rowRes.AuditType, &rowRes.AuditTo, &rowRes.Status, &rowRes.Created, &rowRes.CreatedBy, &rowRes.Updated, &rowRes.UpdatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *eventRepository) AuditEvent(eventID int64, approved bool, byUser string, auditContent string, currentLevel int) (int64, error) {
	eventStatus := 9
	historyStatus := 1
	historyType := "审核通过"
	isActive := 0
	nextLevel := currentLevel
	nextAuditType := 0
	fmt.Println(nextLevel)
	if !approved {
		eventStatus = 3
		isActive = 1
		nextLevel = 1
		row := r.tx.QueryRow(`SELECT audit_type FROM event_audits WHERE event_id = ? AND audit_level = ? AND status > 0 ORDER BY audit_level ASC`, eventID, 1)
		err := row.Scan(&nextAuditType)
		if err != nil {
			return 0, err
		}
	} else {
		row := r.tx.QueryRow(`SELECT audit_level, audit_type FROM event_audits WHERE event_id = ? AND audit_level > ? AND status > 0 ORDER BY audit_level ASC`, eventID, currentLevel)
		err := row.Scan(&nextLevel, &nextAuditType)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				nextLevel = 0
				nextAuditType = 0
			} else {
				return 0, err
			}
		}
		if nextLevel != 0 {
			eventStatus = 2
			isActive = 1
		}
	}
	_, err := r.tx.Exec(`
		Update events SET 
		audit_level = ?,
		audit_type = ?,
		audit_user = ?,
		audit_time = ?,
		audit_content = ?,
		is_active = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, nextLevel, nextAuditType, byUser, time.Now().Format("2006-01-02 15:04:05"), auditContent, isActive, eventStatus, time.Now(), byUser, eventID)
	if err != nil {
		return 0, err
	}
	if approved {
		historyStatus = 1
	} else {
		historyStatus = 2
		historyType = "审核驳回"
	}
	res, err := r.tx.Exec(`
		INSERT INTO event_historys 
		(event_id, history_type, audit_user, audit_content, audit_time, status, created, created_by, updated, updated_by)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, eventID, historyType, byUser, auditContent, time.Now().Format("2006-01-02 15:04:05"), historyStatus, time.Now(), byUser, time.Now(), byUser)
	if err != nil {
		return 0, err
	}
	historyID, err := res.LastInsertId()
	return historyID, err
}

func (r *eventRepository) CheckCheckin(eventID int64, userID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM event_checkins WHERE event_id = ? AND user_id = ? AND status > 0 AND checkin_time > ? AND checkin_time < ? `, eventID, userID, time.Now().Format("2006-01-02")+" 00:00:00", time.Now().Format("2006-01-02")+" 23:59:59")
	err := row.Scan(&res)
	return res, err
}

func (r *eventRepository) doCheckin(eventID int64, info NewCheckin) error {
	_, err := r.tx.Exec(`
		INSERT INTO event_checkins
		(
			event_id,
			user_id,
			user_name,
			checkin_type,
			checkin_time,
			longitude,
			latitude,
			distance,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, eventID, info.UserID, info.User, info.CheckinType, time.Now(), info.Longitude, info.Latitude, info.Distance, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetProjectLocation(projectID, organizationID int64) (float64, float64, int, error) {
	var longitude, latitude float64
	var distance int
	var row *sql.Row
	if organizationID == 0 {
		row = r.tx.QueryRow(`SELECT longitude, latitude, checkin_distance FROM projects WHERE id = ? AND status > 0`, projectID)
	} else {
		row = r.tx.QueryRow(`SELECT longitude, latitude, checkin_distance FROM projects WHERE id = ? AND organization_id = ? AND status > 0`, projectID, organizationID)
	}
	err := row.Scan(&longitude, &latitude, &distance)
	return longitude, latitude, distance, err
}

func (r *eventRepository) CreateEventReview(eventID int64, info EventReviewNew) error {
	_, err := r.tx.Exec(`
		INSERT INTO event_reviews
		(
			event_id,
			result,
			content,
			link,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, eventID, info.Result, info.Content, info.Link, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *eventRepository) UpdateEventDeadline(id int64, deadline, byUser string) error {
	if deadline == "" {
		_, err := r.tx.Exec(`
			Update events SET 
			deadline = null,
			updated = ?,
			updated_by = ? 
			WHERE id = ?
		`, time.Now(), byUser, id)
		return err
	} else {
		_, err := r.tx.Exec(`
			Update events SET 
			deadline = ?,
			updated = ?,
			updated_by = ? 
			WHERE id = ?
		`, deadline, time.Now(), byUser, id)
		return err
	}
}

func (r *eventRepository) GetReviewByID(id int64) (*EventReviewResponse, error) {
	var res EventReviewResponse
	// var row *sql.Row
	// if organizationID != 0 {
	// 	row = r.tx.QueryRow(`
	// 		SELECT
	// 		er.id,
	// 		er.event_id,
	// 		er.result,
	// 		er.content,
	// 		er.link,
	// 		er.status,
	// 		er.created
	// 		FROM event_reviews er
	// 		LEFT JOIN events e
	// 		ON e.id = er.event_Id
	// 		LEFT JOIN projects p
	// 		ON e.project_id = p.id
	// 		WHERE re.id = ?
	// 		AND p.organization_id = ?
	// 		AND er.status > 0
	// 		LIMIT 1`, id, organizationID)
	// } else {
	row := r.tx.QueryRow(`
			SELECT 
			id, 
			event_id, 
			result, 
			content, 
			link, 
			status, 
			created,
			handle_time,
			handle_content,
			handle_user
			FROM event_reviews 
			WHERE id = ? 
			AND status > 0 
			LIMIT 1`, id)
	// }
	err := row.Scan(&res.ID, &res.EventID, &res.Result, &res.Content, &res.Link, &res.Status, &res.Created, &res.HandleTime, &res.HandleContent, &res.HandleUser)
	return &res, err
}

func (r *eventRepository) HandleReview(reviewID int64, status int, byUser string, handleContent string) error {
	_, err := r.tx.Exec(`
		Update event_reviews SET 
		handle_user = ?,
		handle_time = ?,
		handle_content = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, byUser, time.Now().Format("2006-01-02 15:04:05"), handleContent, status, time.Now(), byUser, reviewID)
	return err
}

func (r *eventRepository) SetEventActive(eventID int64) error {
	_, err := r.tx.Exec(`
		UPDATE events SET
		is_active = 1,
		updated = ?
		WHERE id = ?
	`, time.Now(), eventID)
	return err
}

func (r *eventRepository) GetProjectProgress(id int64) (int, int, error) {
	var all, completed int
	row := r.tx.QueryRow(`
		SELECT 
		count(1) 
		FROM events 
		WHERE project_id = ?
		AND status > 0`, id)
	err := row.Scan(&all)
	if err != nil {
		return 0, 0, err
	}
	row = r.tx.QueryRow(`
		SELECT 
		count(1)
		FROM events 
		WHERE project_id = ?
		AND status = 9`, id)
	err = row.Scan(&completed)
	if err != nil {
		return 0, 0, err
	}
	return all, completed, err
}

func (r *eventRepository) UpdateProjectProgress(projectID int64, progress int) error {
	sql := `UPDATE projects set progress = ?,`
	if progress == 100 {
		sql += ` status = 2,`
	}
	sql += ` updated = ? WHERE id = ?`
	_, err := r.tx.Exec(sql, progress, time.Now(), projectID)
	return err
}

func (r *eventRepository) DeleteEventAuditFile(eventID int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update event_audit_files SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE event_id = ?
	`, time.Now(), byUser, eventID)
	return err
}

func (r *eventRepository) CreateEventAuditFile(info EventAuditFile) error {
	_, err := r.tx.Exec(`
		INSERT INTO event_audit_files
		(
			event_id,
			link,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, info.EventID, info.Link, info.Status, time.Now(), info.CreatedBy, time.Now(), info.UpdatedBy)
	return err
}

func (r *eventRepository) CreateEventHistoryFile(info EventHistoryFile) error {
	_, err := r.tx.Exec(`
		INSERT INTO event_history_files
		(
			history_id,
			link,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, info.HistoryID, info.Link, info.Status, time.Now(), info.CreatedBy, time.Now(), info.UpdatedBy)
	return err
}
