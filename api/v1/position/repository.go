package position

import (
	"database/sql"
	"time"
)

type positionRepository struct {
	tx *sql.Tx
}

func NewPositionRepository(transaction *sql.Tx) PositionRepository {
	return &positionRepository{
		tx: transaction,
	}
}

type PositionRepository interface {
	//Position Management
	CreatePosition(PositionNew) (int64, error)
	UpdatePosition(int64, PositionNew) (int64, error)
	GetPositionByID(int64, int64) (*Position, error)
	CheckNameExist(string, int64, int64) (int, error)
}

func (r *positionRepository) CreatePosition(info PositionNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO positions
		(
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *positionRepository) UpdatePosition(id int64, info PositionNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update positions SET 
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *positionRepository) GetPositionByID(id int64, organizationID int64) (*Position, error) {
	var res Position
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, status, created, created_by, updated, updated_by FROM positions WHERE id = ? AND organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, status, created, created_by, updated, updated_by FROM positions WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *positionRepository) CheckNameExist(name string, organizationID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM positions WHERE name = ? AND organization_id = ? AND id != ?  LIMIT 1`, name, organizationID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
