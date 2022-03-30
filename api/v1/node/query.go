package node

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type nodeQuery struct {
	conn *sqlx.DB
}

func NewNodeQuery(connection *sqlx.DB) NodeQuery {
	return &nodeQuery{
		conn: connection,
	}
}

type NodeQuery interface {
	//Node Management
	GetNodeByID(id int64) (*Node, error)
	GetAssignsByNodeID(int64) (*[]NodeAssign, error)
	GetPresByNodeID(int64) (*[]NodePre, error)
	GetNodeCount(NodeFilter, int64) (int, error)
	GetNodeList(NodeFilter, int64) (*[]Node, error)
	//WX
	GetAssigned(int64, int64) ([]int64, error)
	CheckActive(int64) (bool, error)
	GetAssignedNodeByID(int64, string) (*Node, error)
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

func (r *nodeQuery) GetAssigned(userID int64, positionID int64) ([]int64, error) {
	var assigns []int64
	err := r.conn.Select(&assigns, "SELECT node_id FROM node_assigns WHERE ((assign_type = 2 AND assign_to  = ?) OR (assign_type = 1 AND assign_to = ?)) AND status = ?", userID, positionID, 1)
	return assigns, err
}

func (r *nodeQuery) CheckActive(nodeID int64) (bool, error) {
	var activePreCount int
	err := r.conn.Get(&activePreCount, `
		SELECT count(1) from node_pres ep
		LEFT JOIN nodes e
		ON ep.pre_id = e.id 
		WHERE ep.status = 1  
		AND ep.node_id = ?
		AND e.status = 1`, nodeID)
	if err != nil {
		return false, err
	}
	if activePreCount == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *nodeQuery) GetAssignedNodeByID(id int64, status string) (*Node, error) {
	var node Node
	sql := "SELECT * FROM nodes WHERE id = ?"
	if status == "all" {
		sql = sql + " AND status > 0"
	} else {
		sql = sql + " AND status = 1"
	}
	err := r.conn.Get(&node, sql, id)
	if err != nil {
		return nil, err
	}
	return &node, nil
}
