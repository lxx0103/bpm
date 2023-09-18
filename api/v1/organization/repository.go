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
			logo2,
			description,
			contact,
			phone,
			address,
			city,
			type,
			user_limit,
			expiry_date,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Logo, info.Logo2, info.Description, info.Contact, info.Phone, info.Address, info.City, info.Type, info.UserLimit, info.ExpiryDate, info.Status, time.Now(), info.User, time.Now(), info.User)
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
		logo = ?,
		logo2 = ?,
		city = ?,
		type = ?,
		user_limit = ?,
		expiry_date = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Description, info.Contact, info.Phone, info.Address, info.Logo, info.Logo2, info.City, info.Type, info.UserLimit, info.ExpiryDate, info.Status, time.Now(), info.User, id)
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
	row := r.tx.QueryRow(`SELECT id, name, logo, logo2, description, contact, phone, address, city, type, user_limit, IFNULL(expiry_date, "") as expiry_date, status, created, created_by, updated, updated_by FROM organizations WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Logo, &res.Logo2, &res.Description, &res.Contact, &res.Phone, &res.Address, &res.City, &res.Type, &res.UserLimit, &res.ExpiryDate, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
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

func (r *organizationRepository) NewQrcode(path, source, img string) error {
	_, err := r.tx.Exec(`
		INSERT INTO qr_codes
		(
			path,
			img,
			source,
			created,
			created_by
		)
		VALUES (?, ?, ?, now(), "SYSTEM")
	`, path, img, source)
	return err
}

func (r *organizationRepository) CreateOrganizationQrcode(organizationID int64, qrcode OrganizationQrcode, byUser string) error {
	_, err := r.tx.Exec(`
		INSERT INTO organization_qrcodes
		(
			organization_id,
			type,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, organizationID, qrcode.Type, qrcode.Name, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *organizationRepository) DeleteOrganizationQrcode(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update organization_qrcodes SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE organization_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}
