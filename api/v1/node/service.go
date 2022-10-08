package node

import (
	"bpm/core/database"
	"errors"
)

type nodeService struct {
}

func NewNodeService() *nodeService {
	return &nodeService{}
}

func (s *nodeService) GetNodeByID(id int64) (*Node, error) {
	db := database.InitMySQL()
	query := NewNodeQuery(db)
	node, err := query.GetNodeByID(id)
	if err != nil {
		return nil, err
	}
	assigns, err := query.GetAssignsByNodeID(node.ID)
	if err != nil {
		return nil, err
	}
	node.Assign = assigns
	pres, err := query.GetPresByNodeID(node.ID)
	if err != nil {
		return nil, err
	}
	node.PreID = pres
	audits, err := query.GetAuditsByNodeID(node.ID)
	if err != nil {
		return nil, err
	}
	node.Audit = audits
	return node, err
}

func (s *nodeService) NewNode(info NodeNew, organizationID int64) (*Node, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewNodeRepository(tx)
	templateExist, err := repo.CheckTemplateExist(info.TemplateID, organizationID)
	if err != nil {
		return nil, err
	}
	if templateExist == 0 {
		msg := "项目不存在"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, info.TemplateID, 0)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "节点名称重复"
		return nil, errors.New(msg)
	}
	nodeID, err := repo.CreateNode(info)
	if err != nil {
		return nil, err
	}
	node, err := repo.GetNodeByID(nodeID, organizationID)
	if err != nil {
		return nil, err
	}
	if info.AssignType != 3 {
		err = repo.CreateNodeAssign(nodeID, info.AssignType, info.AssignTo, info.User)
		if err != nil {
			return nil, err
		}
	}
	assigns, err := repo.GetAssignsByNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	node.Assign = assigns
	err = repo.CreateNodePre(nodeID, info.PreID, info.User)
	if err != nil {
		return nil, err
	}
	pres, err := repo.GetPresByNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	node.PreID = pres
	if info.NeedAudit == 1 {
		err = repo.CreateNodeAudit(nodeID, info.AuditType, info.AuditTo, info.User)
		if err != nil {
			return nil, err
		}
	}
	audits, err := repo.GetAuditsByNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	node.Audit = audits
	tx.Commit()
	return node, err
}

func (s *nodeService) GetNodeList(filter NodeFilter, organizationID int64) (int, *[]Node, error) {
	db := database.InitMySQL()
	query := NewNodeQuery(db)
	count, err := query.GetNodeCount(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetNodeList(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *nodeService) UpdateNode(nodeID int64, info NodeUpdate, organizationID int64) (*Node, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewNodeRepository(tx)
	oldNode, err := repo.GetNodeByID(nodeID, organizationID)
	if err != nil {
		return nil, err
	}
	if info.Name != "" {
		exist, err := repo.CheckNameExist(info.Name, oldNode.TemplateID, nodeID)
		if err != nil {
			return nil, err
		}
		if exist != 0 {
			msg := "节点名称重复"
			return nil, errors.New(msg)
		}
		oldNode.Name = info.Name
	}
	if info.Assignable != 0 {
		oldNode.Assignable = info.Assignable
	}
	if info.AssignType != 0 {
		oldNode.AssignType = info.AssignType
	}
	if info.NeedAudit != 0 {
		oldNode.NeedAudit = info.NeedAudit
	}
	if info.AuditType != 0 {
		oldNode.AuditType = info.AuditType
	}
	if info.NeedCheckin != 0 {
		oldNode.NeedCheckin = info.NeedCheckin
	}
	if info.Sort != 0 {
		oldNode.Sort = info.Sort
	}
	if info.CanReview != 0 {
		oldNode.CanReview = info.CanReview
	}
	oldNode.JsonData = info.JsonData
	err = repo.UpdateNode(nodeID, *oldNode, info.User)
	if err != nil {
		msg := "更新节点失败"
		return nil, errors.New(msg)
	}
	node, err := repo.GetNodeByID(nodeID, organizationID)
	if err != nil {
		return nil, err
	}
	err = repo.DeleteNodeAssign(nodeID, info.User)
	if err != nil {
		return nil, err
	}
	if info.AssignType != 0 && info.AssignType != 3 {
		err = repo.CreateNodeAssign(nodeID, info.AssignType, info.AssignTo, info.User)
		if err != nil {
			return nil, err
		}
	}
	assigns, err := repo.GetAssignsByNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	node.Assign = assigns
	err = repo.DeleteNodePre(nodeID, info.User)
	if err != nil {
		return nil, err
	}
	if len(info.PreID) != 0 {
		err = repo.CreateNodePre(nodeID, info.PreID, info.User)
		if err != nil {
			return nil, err
		}
	}
	pres, err := repo.GetPresByNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	node.PreID = pres

	err = repo.DeleteNodeAudit(nodeID, info.User)
	if err != nil {
		return nil, err
	}
	if info.AuditType != 0 {
		err = repo.CreateNodeAudit(nodeID, info.AuditType, info.AuditTo, info.User)
		if err != nil {
			return nil, err
		}
	}
	audits, err := repo.GetAuditsByNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	node.Audit = audits
	tx.Commit()
	return node, err
}

func (s *nodeService) DeleteNode(nodeID int64, organizationID int64, user string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewNodeRepository(tx)
	_, err = repo.GetNodeByID(nodeID, organizationID)
	if err != nil {
		return err
	}
	err = repo.DeleteNodeAssign(nodeID, user)
	if err != nil {
		return err
	}
	err = repo.DeleteNodePre(nodeID, user)
	if err != nil {
		return err
	}
	err = repo.DeleteNodeAudit(nodeID, user)
	if err != nil {
		return err
	}
	err = repo.DeleteNode(nodeID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
