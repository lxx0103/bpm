package vendors

import (
	"database/sql"
	"time"
)

type vendorsRepository struct {
	tx *sql.Tx
}

func NewVendorsRepository(transaction *sql.Tx) *vendorsRepository {
	return &vendorsRepository{
		tx: transaction,
	}
}

func (r *vendorsRepository) CreateVendors(info VendorsNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO vendors
		(
			name,
			contact,
			phone,
			address,
			longitude,
			latitude,
			description,
			cover,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Contact, info.Phone, info.Address, info.Longitude, info.Latitude, info.Description, info.Cover, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *vendorsRepository) CheckMaterialExist(id int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM materials WHERE id = ? AND status > 0  LIMIT 1`, id)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *vendorsRepository) CreateVendorsMaterial(vendorsID, materialID int64, byUser string) error {
	_, err := r.tx.Exec(`
		INSERT INTO vendor_materials
		(
			vendor_id,
			material_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, vendorsID, materialID, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *vendorsRepository) CheckBrandExist(id int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM brands WHERE id = ? AND status > 0  LIMIT 1`, id)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *vendorsRepository) CreateVendorsBrand(vendorsID, brandID int64, byUser string) error {
	_, err := r.tx.Exec(`
		INSERT INTO vendor_brands
		(
			vendor_id,
			brand_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, vendorsID, brandID, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *vendorsRepository) CreateVendorsPicture(vendorsID int64, picture, byUser string) error {
	_, err := r.tx.Exec(`
		INSERT INTO vendor_pictures
		(
			vendor_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, vendorsID, picture, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *vendorsRepository) CreateVendorsQrcode(vendorsID int64, qrcode VendorsQrcode, byUser string) error {
	_, err := r.tx.Exec(`
		INSERT INTO vendor_qrcodes
		(
			vendor_id,
			type,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, vendorsID, qrcode.Type, qrcode.Name, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *vendorsRepository) UpdateVendors(id int64, info VendorsNew) error {
	_, err := r.tx.Exec(`
		Update vendors SET 
		name = ?,
		contact = ?,
		phone = ?,
		address = ?,
		longitude = ?,
		latitude = ?,
		description = ?,
		cover = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Contact, info.Phone, info.Address, info.Longitude, info.Latitude, info.Description, info.Cover, time.Now(), info.User, id)
	return err
}

func (r *vendorsRepository) GetVendorsByID(id int64) (*VendorsResponse, error) {
	var res VendorsResponse
	row := r.tx.QueryRow(`SELECT id, name, status FROM vendors WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Status)
	return &res, err
}

func (r *vendorsRepository) CheckVendorsNameExist(name string, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM vendors WHERE name = ? AND status > 0 AND id != ?  LIMIT 1`, name, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *vendorsRepository) CheckVendorsActive(vendorsID int64) (int, error) {
	// var res int
	// row := r.tx.QueryRow(`SELECT count(1) FROM material_vendors WHERE vendor_id = ? AND status > 0 LIMIT 1`, vendorsID)
	// err := row.Scan(&res)
	// if err != nil {
	// 	return 0, err
	// }
	res := 0
	return res, nil
}

func (r *vendorsRepository) DeleteVendors(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendors SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *vendorsRepository) DeleteVendorsMaterial(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendor_materials SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE vendor_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *vendorsRepository) DeleteVendorsBrand(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendor_brands SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE vendor_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *vendorsRepository) DeleteVendorsPicture(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendor_pictures SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE vendor_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *vendorsRepository) DeleteVendorsQrcode(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendor_qrcodes SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE vendor_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}
