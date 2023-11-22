package shortcut

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type shortcutQuery struct {
	conn *sqlx.DB
}

func NewShortcutQuery(connection *sqlx.DB) *shortcutQuery {
	return &shortcutQuery{
		conn: connection,
	}
}

func (r *shortcutQuery) GetShortcutByID(id int64, organizationID int64) (*ShortcutResponse, error) {
	var shortcut ShortcutResponse
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&shortcut, `
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
		AND m.organization_id = ? 
		AND m.status > 0
		`, id, organizationID)
	} else {
		err = r.conn.Get(&shortcut, `		
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
	}
	return &shortcut, err
}

func (r *shortcutQuery) GetShortcutCount(filter ShortcutFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.ShortcutType; v != 0 {
		where, args = append(where, "shortcut_type = ?"), append(args, v)
	}
	if v := filter.Keyword; v != "" {
		where, args = append(where, "content like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM shortcuts 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *shortcutQuery) GetShortcutList(filter ShortcutFilter) (*[]ShortcutResponse, error) {
	where, args := []string{"m.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "m.organization_id = ?"), append(args, v)
	}
	if v := filter.ShortcutType; v != 0 {
		where, args = append(where, "m.shortcut_type = ?"), append(args, v)
	}
	if v := filter.Keyword; v != "" {
		where, args = append(where, "m.content like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var shortcuts []ShortcutResponse
	err := r.conn.Select(&shortcuts, `
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
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &shortcuts, nil
}

func (r *shortcutQuery) GetShortcutTypeList(filter ShortcutTypeFilter) (*[]ShortcutTypeResponse, error) {
	where, args := []string{"m.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "m.organization_id = ?"), append(args, v)
	}
	if v := filter.ParentID; v != 0 {
		where, args = append(where, "m.parent_id = ?"), append(args, v)
	}
	var shortcuts []ShortcutTypeResponse
	err := r.conn.Select(&shortcuts, `
		SELECT 
		m.id,
		m.organization_id, 
		o.name as organization_name,
		m.parent_id,
		m.name,
		m.status
		FROM shortcut_types m
		LEFT JOIN organizations o
		ON m.organization_id = o.id 
		WHERE `+strings.Join(where, " AND "), args...)
	return &shortcuts, err
}

func (r *shortcutQuery) GetShortcutTypeByID(id int64, organizationID int64) (*ShortcutTypeResponse, error) {
	var shortcutType ShortcutTypeResponse
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&shortcutType, `
		SELECT 
		m.id,
		m.organization_id, 
		o.name as organization_name,
		m.parent_id,
		m.name,
		m.status
		FROM shortcut_types m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE m.id = ? 
		AND m.organization_id = ? 
		AND m.status > 0
		`, id, organizationID)
	} else {
		err = r.conn.Get(&shortcutType, `		
		SELECT 
		m.id,
		m.organization_id, 
		o.name as organization_name,
		m.parent_id,
		m.name,
		m.status
		FROM shortcut_types m
		LEFT JOIN organizations o
		ON m.organization_id = o.id
		WHERE m.id = ? 
		AND m.status > 0
		`, id)
	}
	return &shortcutType, err
}
