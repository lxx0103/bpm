package team

import (
	"database/sql"
	"time"
)

type teamRepository struct {
	tx *sql.Tx
}

func NewTeamRepository(transaction *sql.Tx) *teamRepository {
	return &teamRepository{
		tx: transaction,
	}
}

func (r *teamRepository) CreateTeam(info TeamNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO teams
		(
			organization_id,
			name,
			leader,
			phone,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Leader, info.Phone, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *teamRepository) UpdateTeam(id int64, info TeamNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update teams SET 
		name = ?,
		leader = ?,
		phone = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Leader, info.Phone, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *teamRepository) GetTeamByID(id int64, organizationID int64) (*Team, error) {
	var res Team
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, leader, phone, status, created, created_by, updated, updated_by FROM teams WHERE id = ? AND organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, leader, phone, status, created, created_by, updated, updated_by FROM teams WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Leader, &res.Phone, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *teamRepository) CheckNameExist(name string, organizationID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM teams WHERE name = ? AND organization_id = ? AND id != ?  LIMIT 1`, name, organizationID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
