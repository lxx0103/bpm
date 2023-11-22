package shortcut

import (
	"database/sql"
	"time"
)

type shortcutRepository struct {
	tx *sql.Tx
}

func NewShortcutRepository(transaction *sql.Tx) *shortcutRepository {
	return &shortcutRepository{
		tx: transaction,
	}
}

func (r *shortcutRepository) CreateShortcut(info ShortcutNew) (int64, error) {
	res, err := r.tx.Exec(`
		INSERT INTO shortcuts
		(
			organization_id,
			shortcut_type,
			shortcut_type_name,
			content,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ShortcutType, info.ShortcutTypeName, info.Content, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	shortcutID, err := res.LastInsertId()
	return shortcutID, err
}

func (r *shortcutRepository) UpdateShortcut(id int64, info ShortcutUpdate) error {
	_, err := r.tx.Exec(`
		Update shortcuts SET 
		content = ?,
		shortcut_type_name = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Content, info.ShortcutTypeName, time.Now(), info.User, id)
	return err
}

func (r *shortcutRepository) GetShortcutByID(id int64) (*ShortcutResponse, error) {
	var res ShortcutResponse
	row := r.tx.QueryRow(`
		SELECT 
		m.id,
		m.organization_id, 
		o.name as organization_name,
		m.shortcut_type,
		m.shortcut_type_name,
		m.content, 
		m.status
		FROM shortcuts m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE m.id = ? 
		AND m.status > 0
	`, id)

	err := row.Scan(&res.ID, &res.OrganizationID, &res.OrganizationName, &res.ShortcutType, &res.ShortcutTypeName, &res.Content, &res.Status)
	return &res, err
}

func (r *shortcutRepository) DeleteShortcut(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update shortcuts SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *shortcutRepository) GetShortcutTypeByID(id int64) (*ShortcutTypeResponse, error) {
	var res ShortcutTypeResponse
	row := r.tx.QueryRow(`
		SELECT 
		m.id,
		m.organization_id, 
		o.name as organization_name,
		m.name,
		m.parent_id,
		m.status
		FROM shortcut_types m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE m.id = ? 
		AND m.status > 0
	`, id)

	err := row.Scan(&res.ID, &res.OrganizationID, &res.OrganizationName, &res.Name, &res.ParentID, &res.Status)
	return &res, err
}

func (r *shortcutRepository) CreateShortcutType(info ShortcutTypeNew) error {
	_, err := r.tx.Exec(`
		INSERT INTO shortcut_types
		(
			organization_id,
			parent_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ParentID, info.Name, 1, time.Now(), info.User, time.Now(), info.User)
	return err
}

func (r *shortcutRepository) UpdateShortcutType(id int64, info ShortcutTypeUpdate) error {
	_, err := r.tx.Exec(`
		Update shortcut_types SET 
		parent_id = ?,
		name = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.ParentID, info.Name, time.Now(), info.User, id)
	return err
}

func (r *shortcutRepository) DeleteShortcutType(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update shortcut_types SET 
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), byUser, id)
	return err
}

func (r *shortcutRepository) GetChildCount(typeID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`
		SELECT count(1) 
		FROM shortcut_types
		WHERE parent_id = ? 
		AND status > 0 
		`, typeID)
	err := row.Scan(&res)
	return res, err
}
