package client

import (
	"database/sql"
	"time"
)

type clientRepository struct {
	tx *sql.Tx
}

func NewClientRepository(transaction *sql.Tx) ClientRepository {
	return &clientRepository{
		tx: transaction,
	}
}

type ClientRepository interface {
	//Client Management
	CreateClient(ClientNew, int64) (int64, error)
	UpdateClient(int64, ClientNew) (int64, error)
	GetClientByID(int64, int64) (*Client, error)
}

func (r *clientRepository) CreateClient(info ClientNew, organizationID int64) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO clients
		(
			organization_id,
			name,
			phone,
			address,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, organizationID, info.Name, info.Phone, info.Address, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *clientRepository) UpdateClient(id int64, info ClientNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update clients SET 
		name = ?,
		phone = ?,
		address = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Phone, info.Address, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *clientRepository) GetClientByID(id int64, organizationID int64) (*Client, error) {
	var res Client
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, phone, address, status, created, created_by, updated, updated_by FROM clients WHERE id = ? AND organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, phone, address, status, created, created_by, updated, updated_by FROM clients WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Phone, &res.Address, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
