package node

import (
	"database/sql"
	"errors"
	"time"
)

type nodeRepository struct {
	tx *sql.Tx
}

func NewNodeRepository(transaction *sql.Tx) *nodeRepository {
	return &nodeRepository{
		tx: transaction,
	}
}

func (r *nodeRepository) CreateNode(info NodeNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO nodes
		(
			template_id,
			name,
			assignable,
			assign_type,
			need_audit,
			audit_type,
			need_checkin,
			sort,
			status,
			json_data,
			can_review,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.TemplateID, info.Name, info.Assignable, info.AssignType, info.NeedAudit, info.AuditType, info.NeedCheckin, info.Sort, 1, "{}", info.CanReview, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *nodeRepository) UpdateNode(id int64, info Node, byUser string) error {
	_, err := r.tx.Exec(`
		Update nodes SET 
		name = ?,
		assignable = ?,
		assign_type = ?,
		need_audit = ?,
		audit_type = ?,
		need_checkin = ?,
		sort = ?,
		can_review = ?,
		json_data = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Assignable, info.AssignType, info.NeedAudit, info.AuditType, info.NeedCheckin, info.Sort, info.CanReview, info.JsonData, time.Now(), byUser, id)
	return err
}

func (r *nodeRepository) CreateNodeAssign(nodeID int64, assignType int, assignTo []int64, user string) error {
	for i := 0; i < len(assignTo); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM node_assigns WHERE node_id = ? AND assign_type = ? AND assign_to = ? AND status > 0  LIMIT 1`, nodeID, assignType, assignTo[i])
		err := row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist != 0 {
			msg := "指派对象有重复"
			return errors.New(msg)
		}
		_, err = r.tx.Exec(`
			INSERT INTO node_assigns
			(
				node_id,
				assign_type,
				assign_to,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, nodeID, assignType, assignTo[i], 1, time.Now(), user, time.Now(), user)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *nodeRepository) DeleteNodeAssign(node_id int64, user string) error {
	_, err := r.tx.Exec(`
		Update node_assigns SET
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE node_id = ?
	`, -1, time.Now(), user, node_id)
	return err
}

func (r *nodeRepository) GetAssignsByNodeID(nodeID int64) (*[]NodeAssign, error) {
	var res []NodeAssign
	rows, err := r.tx.Query(`SELECT id, node_id, assign_type, assign_to, status, created, created_by, updated, updated_by FROM node_assigns WHERE node_id = ? AND status > 0 `, nodeID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes NodeAssign
		err = rows.Scan(&rowRes.ID, &rowRes.NodeID, &rowRes.AssignType, &rowRes.AssignTo, &rowRes.Status, &rowRes.Created, &rowRes.CreatedBy, &rowRes.Updated, &rowRes.UpdatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *nodeRepository) GetNodeByID(id int64, organizationID int64) (*Node, error) {
	var res Node
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT e.id, e.template_id, e.name, e.assignable, e.assign_type, e.need_audit, e.audit_type, e.need_checkin, e.sort, e.can_review, e.status, e.json_data, e.created, e.created_by, e.updated, e.updated_by FROM nodes e LEFT JOIN templates p ON e.template_id = p.id  WHERE e.id = ? AND p.organization_id = ? AND e.status > 0 LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, template_id, name, assignable, assign_type, need_audit, audit_type, need_checkin, sort, can_review, status, json_data, created, created_by, updated, updated_by FROM nodes WHERE id = ? AND status > 0 LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.TemplateID, &res.Name, &res.Assignable, &res.AssignType, &res.NeedAudit, &res.AuditType, &res.NeedCheckin, &res.Sort, &res.CanReview, &res.Status, &res.JsonData, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *nodeRepository) CheckTemplateExist(templateID int64, organizationID int64) (int, error) {
	var res int
	var row *sql.Row
	if organizationID == 0 {
		row = r.tx.QueryRow(`SELECT count(1) FROM templates WHERE id = ? AND status > 0  LIMIT 1`, templateID)
	} else {
		row = r.tx.QueryRow(`SELECT count(1) FROM templates WHERE id = ? AND organization_id = ? AND status > 0  LIMIT 1`, templateID, organizationID)
	}
	err := row.Scan(&res)
	return res, err
}

func (r *nodeRepository) CheckNameExist(name string, templateID int64, selfID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM nodes WHERE name = ? AND template_id = ? AND id != ? AND status > 0  LIMIT 1`, name, templateID, selfID)
	err := row.Scan(&res)
	return res, err
}

func (r *nodeRepository) CreateNodePre(nodeID int64, preIDs []int64, user string) error {
	for i := 0; i < len(preIDs); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM node_pres WHERE node_id = ? AND pre_id = ? AND status > 0  LIMIT 1`, nodeID, preIDs[i])
		err := row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist != 0 {
			msg := "前置节点有重复"
			return errors.New(msg)
		}
		_, err = r.tx.Exec(`
			INSERT INTO node_pres
			(
				node_id,
				pre_id,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, nodeID, preIDs[i], 1, time.Now(), user, time.Now(), user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *nodeRepository) DeleteNodePre(node_id int64, user string) error {
	_, err := r.tx.Exec(`
		Update node_pres SET
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE node_id = ?
	`, -1, time.Now(), user, node_id)
	return err
}

func (r *nodeRepository) GetPresByNodeID(nodeID int64) (*[]NodePre, error) {
	var res []NodePre
	rows, err := r.tx.Query(`SELECT id, node_id, pre_id, status, created, created_by, updated, updated_by FROM node_pres WHERE node_id = ? AND status > 0`, nodeID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes NodePre
		err = rows.Scan(&rowRes.ID, &rowRes.NodeID, &rowRes.PreID, &rowRes.Status, &rowRes.Created, &rowRes.CreatedBy, &rowRes.Updated, &rowRes.UpdatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *nodeRepository) DeleteNode(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update node_pres SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE node_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		Update nodes SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *nodeRepository) GetNodesByTemplateID(templateID int64) (*[]Node, error) {
	var res []Node
	rows, err := r.tx.Query(`SELECT id, template_id, name, assign_type, assignable, need_audit, audit_type, need_checkin, sort, can_review FROM nodes  WHERE template_id = ? AND status > 0`, templateID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes Node
		err = rows.Scan(&rowRes.ID, &rowRes.TemplateID, &rowRes.Name, &rowRes.AssignType, &rowRes.Assignable, &rowRes.NeedAudit, &rowRes.AuditType, &rowRes.NeedCheckin, &rowRes.Sort, &rowRes.CanReview)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}

func (r *nodeRepository) CreateNodeAudit(nodeID int64, auditLevel, auditType int, auditTo []int64, user string) error {
	for i := 0; i < len(auditTo); i++ {
		var exist int
		row := r.tx.QueryRow(`SELECT count(1) FROM node_audits WHERE node_id = ? AND audit_level = ? AND audit_type = ? AND audit_to = ? AND status > 0  LIMIT 1`, nodeID, auditLevel, auditType, auditTo[i])
		err := row.Scan(&exist)
		if err != nil {
			return err
		}
		if exist != 0 {
			msg := "审核对象有重复"
			return errors.New(msg)
		}
		_, err = r.tx.Exec(`
			INSERT INTO node_audits
			(
				node_id,
				audit_level,
				audit_type,
				audit_to,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, nodeID, auditLevel, auditType, auditTo[i], 1, time.Now(), user, time.Now(), user)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *nodeRepository) DeleteNodeAudit(node_id int64, user string) error {
	_, err := r.tx.Exec(`
		Update node_audits SET
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE node_id = ?
	`, -1, time.Now(), user, node_id)
	return err
}

func (r *nodeRepository) GetAuditsByNodeID(nodeID int64) (*[]NodeAudit, error) {
	var res []NodeAudit
	rows, err := r.tx.Query(`SELECT id, node_id, audit_level, audit_type, audit_to, status, created, created_by, updated, updated_by FROM node_audits WHERE node_id = ? AND status > 0 `, nodeID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowRes NodeAudit
		err = rows.Scan(&rowRes.ID, &rowRes.NodeID, &rowRes.AuditLevel, &rowRes.AuditType, &rowRes.AuditTo, &rowRes.Status, &rowRes.Created, &rowRes.CreatedBy, &rowRes.Updated, &rowRes.UpdatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, rowRes)
	}
	return &res, nil
}
