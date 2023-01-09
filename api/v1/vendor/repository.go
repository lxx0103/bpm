package vendor

import (
	"database/sql"
	"time"
)

type vendorRepository struct {
	tx *sql.Tx
}

func NewVendorRepository(transaction *sql.Tx) *vendorRepository {
	return &vendorRepository{
		tx: transaction,
	}
}

func (r *vendorRepository) CreateVendor(info VendorNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO vendors
		(
			name,
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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Phone, info.Address, info.Longitude, info.Latitude, info.Description, info.Cover, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *vendorRepository) CheckMaterialExist(id int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM materials WHERE id = ? AND status > 0  LIMIT 1`, id)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *vendorRepository) CreateVendorMaterial(vendorID, materialID int64, byUser string) error {
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
	`, vendorID, materialID, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *vendorRepository) CheckBrandExist(id int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM brands WHERE id = ? AND status > 0  LIMIT 1`, id)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *vendorRepository) CreateVendorBrand(vendorID, brandID int64, byUser string) error {
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
	`, vendorID, brandID, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *vendorRepository) CreateVendorPicture(vendorID int64, picture, byUser string) error {
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
	`, vendorID, picture, 1, time.Now(), byUser, time.Now(), byUser)
	return err
}

func (r *vendorRepository) UpdateVendor(id int64, info VendorNew) error {
	_, err := r.tx.Exec(`
		Update vendors SET 
		name = ?,
		phone = ?,
		address = ?,
		longitude = ?,
		latitude = ?,
		description = ?,
		cover = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Phone, info.Address, info.Longitude, info.Latitude, info.Description, info.Cover, time.Now(), info.User, id)
	return err
}

func (r *vendorRepository) GetVendorByID(id int64) (*VendorResponse, error) {
	var res VendorResponse
	row := r.tx.QueryRow(`SELECT id, name, status FROM vendors WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Status)
	return &res, err
}

func (r *vendorRepository) CheckVendorNameExist(name string, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM vendors WHERE name = ? AND status > 0 AND id != ?  LIMIT 1`, name, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *vendorRepository) CheckVendorActive(vendorID int64) (int, error) {
	// var res int
	// row := r.tx.QueryRow(`SELECT count(1) FROM material_vendors WHERE vendor_id = ? AND status > 0 LIMIT 1`, vendorID)
	// err := row.Scan(&res)
	// if err != nil {
	// 	return 0, err
	// }
	res := 0
	return res, nil
}

func (r *vendorRepository) DeleteVendor(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendors SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *vendorRepository) DeleteVendorMaterial(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendor_materials SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE vendor_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *vendorRepository) DeleteVendorBrand(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendor_brands SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE vendor_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *vendorRepository) DeleteVendorPicture(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update vendor_pictures SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE vendor_id = ?
	`, -1, time.Now(), byUser, id)
	return err
}
