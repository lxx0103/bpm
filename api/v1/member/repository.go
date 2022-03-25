package member

import (
	"database/sql"
	"errors"
	"time"
)

type memberRepository struct {
	tx *sql.Tx
}

func NewMemberRepository(transaction *sql.Tx) MemberRepository {
	return &memberRepository{
		tx: transaction,
	}
}

type MemberRepository interface {
	//Member Management
	CreateProjectMember(int64, []int64, int64, string) error
	DeleteProjectMember(int64, string) error
	GetMembersByProjectID(int64) (*[]MemberResponse, error)
	CheckProjectExist(int64, int64) (int, error)
}

func (r *memberRepository) CreateProjectMember(projectID int64, userID []int64, organizationID int64, user string) error {
	for i := 0; i < len(userID); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM project_members WHERE project_id = ? AND user_id = ? AND status = 1  LIMIT 1`, projectID, userID[i])
		err := row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist != 0 {
			msg := "指派对象有重复"
			return errors.New(msg)
		}
		if organizationID == 0 {
			row = r.tx.QueryRow(`SELECT count(1) FROM users WHERE  id = ? AND status = 1  LIMIT 1`, userID[i])

		} else {
			row = r.tx.QueryRow(`SELECT count(1) FROM users WHERE organization_id = ? AND id = ? AND status = 1  LIMIT 1`, organizationID, userID[i])
		}
		err = row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist == 0 {
			msg := "指派对象不存在"
			return errors.New(msg)
		}
		_, err = r.tx.Exec(`
			INSERT INTO project_members
			(
				project_id,
				user_id,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, projectID, userID[i], 1, time.Now(), user, time.Now(), user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *memberRepository) DeleteProjectMember(projectID int64, user string) error {
	_, err := r.tx.Exec(`
		Update project_members SET
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE project_id = ?
	`, 2, time.Now(), user, projectID)
	if err != nil {
		return err
	}
	return nil
}

func (r *memberRepository) CheckProjectExist(projectID int64, organizationID int64) (int, error) {
	var res int
	var row *sql.Row
	if organizationID == 0 {
		row = r.tx.QueryRow(`SELECT count(1) FROM projects WHERE id = ?  LIMIT 1`, projectID)
	} else {
		row = r.tx.QueryRow(`SELECT count(1) FROM projects WHERE id = ? AND organization_id = ?  LIMIT 1`, projectID, organizationID)
	}
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *memberRepository) CheckNameExist(name string, projectID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM members WHERE name = ? AND project_id = ? AND id != ?  LIMIT 1`, name, projectID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *memberRepository) GetMembersByProjectID(projectID int64) (*[]MemberResponse, error) {
	var res []MemberResponse
	rows, err := r.tx.Query(`
		SELECT u.id, u.name 
		FROM project_members m 
		LEFT JOIN users u 
		ON m.user_id = u.id 
		WHERE m.project_id = ? 
		AND m.status = ? 
		`, projectID, 1)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes MemberResponse
		err = rows.Scan(&rowRes.UserID, &rowRes.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}
