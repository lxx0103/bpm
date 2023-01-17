package common

import (
	"database/sql"
	"time"
)

type commonRepository struct {
	tx *sql.Tx
}

func NewCommonRepository(transaction *sql.Tx) *commonRepository {
	return &commonRepository{
		tx: transaction,
	}
}

func (r *commonRepository) CreateBrand(info BrandNew) error {
	_, err := r.tx.Exec(`
		INSERT INTO brands
		(
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`, info.Name, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *commonRepository) UpdateBrand(id int64, info BrandNew) error {
	_, err := r.tx.Exec(`
		Update brands SET 
		name = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, time.Now(), info.User, id)
	return err
}

func (r *commonRepository) GetBrandByID(id int64) (*BrandResponse, error) {
	var res BrandResponse
	row := r.tx.QueryRow(`SELECT id, name, status FROM brands WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Status)
	return &res, err
}

func (r *commonRepository) CheckBrandNameExist(name string, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM brands WHERE name = ? AND status > 0 AND id != ?  LIMIT 1`, name, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *commonRepository) CheckBrandActive(brandID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM vendor_brands WHERE brand_id = ? AND status > 0 LIMIT 1`, brandID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *commonRepository) DeleteBrand(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update brands SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *commonRepository) CreateMaterial(info MaterialNew) error {
	_, err := r.tx.Exec(`
		INSERT INTO materials
		(
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`, info.Name, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *commonRepository) UpdateMaterial(id int64, info MaterialNew) error {
	_, err := r.tx.Exec(`
		Update materials SET 
		name = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, time.Now(), info.User, id)
	return err
}

func (r *commonRepository) GetMaterialByID(id int64) (*MaterialResponse, error) {
	var res MaterialResponse
	row := r.tx.QueryRow(`SELECT id, name, status FROM materials WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Status)
	return &res, err
}

func (r *commonRepository) CheckMaterialNameExist(name string, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM materials WHERE name = ? AND status > 0 AND id != ?  LIMIT 1`, name, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *commonRepository) CheckMaterialActive(materialID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM vendor_materials WHERE material_id = ? AND status > 0 LIMIT 1`, materialID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *commonRepository) DeleteMaterial(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update materials SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *commonRepository) CreateBanner(info BannerNew) error {
	_, err := r.tx.Exec(`
		INSERT INTO banners
		(
			name,
			picture,
			priority,
			url,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Picture, info.Priority, info.Url, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *commonRepository) UpdateBanner(id int64, info BannerNew) error {
	_, err := r.tx.Exec(`
		Update banners SET 
		name = ?,
		picture = ?,
		priority = ?,
		url = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Picture, info.Priority, info.Url, time.Now(), info.User, id)
	return err
}

func (r *commonRepository) GetBannerByID(id int64) (*BannerResponse, error) {
	var res BannerResponse
	row := r.tx.QueryRow(`SELECT id, name, picture, priority, url, status FROM banners WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Picture, &res.Priority, &res.Url, &res.Status)
	return &res, err
}

func (r *commonRepository) DeleteBanner(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update banners SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}
