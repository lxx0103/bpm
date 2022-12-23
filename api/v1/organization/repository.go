package organization

import (
	"database/sql"
	"time"
)

type organizationRepository struct {
	tx *sql.Tx
}

func NewOrganizationRepository(transaction *sql.Tx) *organizationRepository {
	return &organizationRepository{
		tx: transaction,
	}
}

func (r *organizationRepository) CreateOrganization(info OrganizationNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO organizations
		(
			name,
			logo,
			description,
			contact,
			phone,
			address,
			city,
			type,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Logo, info.Description, info.Contact, info.Phone, info.Address, info.City, info.Type, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *organizationRepository) UpdateOrganization(id int64, info OrganizationNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update organizations SET 
		name = ?,
		description = ?, 
		contact = ?,
		phone = ?,
		address = ?,
		city = ?,
		type = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Description, info.Contact, info.Phone, info.Address, info.City, info.Type, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *organizationRepository) GetOrganizationByID(id int64) (*Organization, error) {
	var res Organization
	row := r.tx.QueryRow(`SELECT id, name, logo, description, contact, phone, address, city, type, status, created, created_by, updated, updated_by FROM organizations WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Logo, &res.Description, &res.Contact, &res.Phone, &res.Address, &res.City, &res.Type, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *organizationRepository) NewAccessToken(code, token string) error {
	_, err := r.tx.Exec(`
		INSERT INTO wx_access_token
		(
			code,
			access_token,
			expires_in
		)
		VALUES (?, ?, DATE_ADD(now(), INTERVAL 2 HOUR))
	`, code, token)
	return err
}

func (r *organizationRepository) NewQrcode(path, img string) error {
	_, err := r.tx.Exec(`
		INSERT INTO qr_codes
		(
			path,
			img,
			created,
			created_by
		)
		VALUES (?, ?, now(), "SYSTEM")
	`, path, img)
	return err
}
