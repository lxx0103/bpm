package element

import (
	"database/sql"
	"time"
)

type elementRepository struct {
	tx *sql.Tx
}

func NewElementRepository(transaction *sql.Tx) ElementRepository {
	return &elementRepository{
		tx: transaction,
	}
}

type ElementRepository interface {
	//Element Management
	CreateElement(info ElementNew) (int64, error)
	UpdateElement(int64, Element, string) error
	GetElementByID(int64) (*Element, error)
	DeleteElement(int64, string) error
	CheckNodeExist(int64, int64) (int, error)
	CheckNameExist(string, int64, int64) (int, error)
	CheckSortExist(int, int64, int64) (int, error)
	GetElementsByNodeID(int64) (*[]Element, error)
}

func (r *elementRepository) CreateElement(info ElementNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO elements
		(
			node_id,
			sort,
			element_type,
			name,
			default_value,
			required,
			patterns,
			json_data,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.NodeID, info.Sort, info.Type, info.Name, info.DefaultValue, info.Required, info.Patterns, info.JsonData, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *elementRepository) UpdateElement(id int64, info Element, byUser string) error {
	_, err := r.tx.Exec(`
		Update elements SET 
		element_type = ?,
		sort = ?,
		name = ?,
		default_value = ?,
		patterns = ?,
		required = ?,
		json_data = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.ElementType, info.Sort, info.Name, info.DefaultValue, info.Patterns, info.Required, info.JsonData, time.Now(), byUser, id)
	return err
}

func (r *elementRepository) GetElementByID(id int64) (*Element, error) {
	var res Element
	row := r.tx.QueryRow(`SELECT id, node_id, sort, element_type, name, value, default_value, required, patterns, json_data, status, created, created_by, updated, updated_by FROM elements WHERE status > 0 AND id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.NodeID, &res.Sort, &res.ElementType, &res.Name, &res.Value, &res.DefaultValue, &res.Required, &res.Patterns, &res.JsonData, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	return &res, err
}

func (r *elementRepository) DeleteElement(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update elements SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *elementRepository) CheckNodeExist(nodeID int64, organizationID int64) (int, error) {
	var res int
	var row *sql.Row
	if organizationID == 0 {
		row = r.tx.QueryRow(`SELECT count(1) FROM nodes n LEFT JOIN templates t ON n.template_id = t.id WHERE n.id = ? AND n.status > 0 AND t.status > 0 LIMIT 1`, nodeID)
	} else {
		row = r.tx.QueryRow(`SELECT count(1) FROM nodes n LEFT JOIN templates t ON n.template_id = t.id WHERE n.id = ? AND n.status > 0 AND t.status > 0 AND t.organization_id = ? LIMIT 1`, nodeID, organizationID)
	}
	err := row.Scan(&res)
	return res, err
}

func (r *elementRepository) CheckNameExist(name string, nodeID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM elements WHERE name = ? AND node_id = ? AND id != ? AND status > 0  LIMIT 1`, name, nodeID, selfID)
	err := row.Scan(&res)
	return res, err
}

func (r *elementRepository) CheckSortExist(sorting int, nodeID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM elements WHERE sort = ? AND node_id = ? AND id != ? AND status > 0  LIMIT 1`, sorting, nodeID, selfID)
	err := row.Scan(&res)
	return res, err
}

func (r *elementRepository) GetElementsByNodeID(nodeID int64) (*[]Element, error) {
	var res []Element
	rows, err := r.tx.Query(`SELECT id, node_id, sort, element_type, name, value, default_value, required, patterns FROM elements WHERE node_id = ? AND status > 0`, nodeID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes Element
		err = rows.Scan(&rowRes.ID, &rowRes.NodeID, &rowRes.Sort, &rowRes.ElementType, &rowRes.Name, &rowRes.Value, &rowRes.DefaultValue, &rowRes.Required, &rowRes.Patterns)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}
