package meeting

import (
	"database/sql"
	"time"
)

type meetingRepository struct {
	tx *sql.Tx
}

func NewMeetingRepository(transaction *sql.Tx) *meetingRepository {
	return &meetingRepository{
		tx: transaction,
	}
}

func (r *meetingRepository) CreateMeeting(info MeetingNew) error {
	_, err := r.tx.Exec(`
		INSERT INTO meetings
		(
			organization_id,
			name,
			date,
			content,
			file,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Date, info.Content, info.File, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *meetingRepository) UpdateMeeting(id int64, info MeetingNew) error {
	_, err := r.tx.Exec(`
		Update meetings SET 
		name = ?,
		date = ?,
		content = ?,
		file = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Date, info.Content, info.File, time.Now(), info.User, id)
	return err
}

func (r *meetingRepository) GetMeetingByID(id int64) (*MeetingResponse, error) {
	var res MeetingResponse
	row := r.tx.QueryRow(`
		SELECT m.id, m.name, m.status, m.organization_id, o.name as organization_name, m.date, m.content, m.file
		FROM meetings m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE m.id = ? 
		AND m.status > 0
	`, id)

	err := row.Scan(&res.ID, &res.Name, &res.Status, &res.OrganizationID, &res.OrganizationName, &res.Date, &res.Content, &res.File)
	return &res, err
}

func (r *meetingRepository) CheckNameExist(name string, organizationID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM meetings WHERE name = ? AND organization_id = ? AND id != ? AND status > 0 LIMIT 1`, name, organizationID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *meetingRepository) DeleteMeeting(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update meetings SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}
