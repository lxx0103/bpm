package example

import (
	"database/sql"
	"time"
)

type exampleRepository struct {
	tx *sql.Tx
}

func NewExampleRepository(transaction *sql.Tx) *exampleRepository {
	return &exampleRepository{
		tx: transaction,
	}
}

func (r *exampleRepository) CreateExample(info ExampleNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO examples
		(
			organization_id,
			name,
			cover,
			style,
			type,
			room,
			notes,
			description,
			example_type,
			finder_user_name,
			feed_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Cover, info.Style, info.Type, info.Room, info.Notes, info.Description, info.ExampleType, info.FinderUserName, info.FeedID, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *exampleRepository) UpdateExample(id int64, info ExampleNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update examples SET 
		name = ?,
		cover = ?,
		organization_id = ?,
		style = ?,
		type = ?,
		room = ?,
		notes = ?,
		description = ?,
		example_type = ?,
		finder_user_name = ?,
		feed_id = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Cover, info.OrganizationID, info.Style, info.Type, info.Room, info.Notes, info.Description, info.ExampleType, info.FinderUserName, info.FeedID, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *exampleRepository) GetExampleByID(id int64, organizationID int64) (*Example, error) {
	var res Example
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, cover, style, type, room, notes, description, example_type, finder_user_name, feed_id, status, created, created_by, updated, updated_by FROM examples WHERE id = ? AND organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, cover, style, type, room, notes, description, example_type, finder_user_name, feed_id, status, created, created_by, updated, updated_by FROM examples WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Cover, &res.Style, &res.Type, &res.Room, &res.Notes, &res.Description, &res.ExampleType, &res.FinderUserName, &res.FeedID, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *exampleRepository) CheckNameExist(name string, organizationID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM examples WHERE name = ? AND organization_id = ? AND id != ?  LIMIT 1`, name, organizationID, selfID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *exampleRepository) CreateExampleMaterial(info ExampleMaterialNew, exampleID int64) error {
	_, err := r.tx.Exec(`
		INSERT INTO example_materials
		(
			example_id,
			vendor_id,
			material_id,
			brand_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, exampleID, info.VendorID, info.MaterialID, info.BrandID, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *exampleRepository) GetExampleMaterialByID(id int64) (*ExampleMaterialResponse, error) {
	var res ExampleMaterialResponse
	row := r.tx.QueryRow(`	
		SELECT em.id, 
		em.example_id, IFNULL(e.name, "") as example_name,
		em.material_id, IFNULL(m.name, "") as material_name,
		em.vendor_id, IFNULL(v.name, "") as vendor_name,
		em.brand_id, IFNULL(b.name, "") as brand_name,
		em.status
		FROM example_materials em 
		LEFT JOIN examples e ON em.example_id = e.id 
		LEFT JOIN materials m ON em.material_id = m.id 
		LEFT JOIN vendors v ON em.vendor_id = v.id 
		LEFT JOIN brands b ON em.brand_id = b.id 
		WHERE em.id = ?
		AND em.status > 0
	`, id)

	err := row.Scan(&res.ID, &res.ExampleID, &res.ExampleName, &res.MaterialID, &res.MaterialName, &res.VendorID, &res.VendorName, &res.BrandID, &res.BrandName, &res.Status)
	return &res, err
}

func (r *exampleRepository) UpdateExampleMaterial(info ExampleMaterialNew, id int64) error {
	_, err := r.tx.Exec(`
		Update example_materials SET 
		material_id = ?,
		vendor_id = ?,
		brand_id = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.MaterialID, info.VendorID, info.BrandID, time.Now(), info.User, id)
	return err
}

func (r *exampleRepository) DeleteExampleMaterial(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update example_materials SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}
