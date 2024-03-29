package node

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type nodeQuery struct {
	conn *sqlx.DB
}

func NewNodeQuery(connection *sqlx.DB) *nodeQuery {
	return &nodeQuery{
		conn: connection,
	}
}

func (r *nodeQuery) GetNodeByID(id int64) (*Node, error) {
	var node Node
	err := r.conn.Get(&node, "SELECT * FROM nodes WHERE id = ? AND status > 0 ", id)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *nodeQuery) GetNodeCount(filter NodeFilter, organizationID int64) (int, error) {
	where, args := []string{"e.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "e.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.TemplateID; v != 0 {
		where, args = append(where, "e.template_id = ?"), append(args, v)
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM nodes e
		LEFT JOIN templates p
		ON e.template_id = p.id
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *nodeQuery) GetNodeList(filter NodeFilter, organizationID int64) (*[]Node, error) {
	where, args := []string{"e.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "e.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.TemplateID; v != 0 {
		where, args = append(where, "e.template_id = ?"), append(args, v)
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var nodes []Node
	err := r.conn.Select(&nodes, `
		SELECT e.* 
		FROM nodes e
		LEFT JOIN templates p
		ON e.template_id = p.id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &nodes, nil
}

func (r *nodeQuery) GetAssignsByNodeID(nodeID int64) (*[]NodeAssign, error) {
	var assigns []NodeAssign
	err := r.conn.Select(&assigns, "SELECT * FROM node_assigns WHERE node_id = ? AND status = ?", nodeID, 1)
	if err != nil {
		return nil, err
	}
	return &assigns, nil
}

func (r *nodeQuery) GetPresByNodeID(nodeID int64) (*[]NodePre, error) {
	var pres []NodePre
	err := r.conn.Select(&pres, "SELECT * FROM node_pres WHERE node_id = ? AND status = ?", nodeID, 1)
	if err != nil {
		return nil, err
	}
	return &pres, nil
}

func (r *nodeQuery) GetAuditsByNodeID(nodeID int64) (*[]NodeAudit, error) {
	var audits []NodeAudit
	err := r.conn.Select(&audits, "SELECT * FROM node_audits WHERE node_id = ? AND status = ?", nodeID, 1)
	if err != nil {
		return nil, err
	}
	return &audits, nil
}
