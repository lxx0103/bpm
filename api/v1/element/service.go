package element

import (
	"bpm/core/database"
	"errors"
)

type elementService struct {
}

func NewElementService() ElementService {
	return &elementService{}
}

// ElementService represents a service for managing elements.
type ElementService interface {
	//Element Management
	GetElementByID(int64) (*Element, error)
	NewElement(ElementNew, int64) (*Element, error)
	GetElementList(ElementFilter) (int, *[]Element, error)
	UpdateElement(int64, ElementUpdate, int64) (*Element, error)
	DeleteElement(int64, string) error
}

func (s *elementService) GetElementByID(id int64) (*Element, error) {
	db := database.InitMySQL()
	query := NewElementQuery(db)
	element, err := query.GetElementByID(id)
	return element, err
}

func (s *elementService) NewElement(info ElementNew, organizationID int64) (*Element, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewElementRepository(tx)
	nodeExist, err := repo.CheckNodeExist(info.NodeID, organizationID)
	if err != nil {
		return nil, err
	}
	if nodeExist == 0 {
		msg := "节点不存在"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, info.NodeID, 0)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "元素名称重复"
		return nil, errors.New(msg)
	}
	exist, err = repo.CheckSortExist(info.Sort, info.NodeID, 0)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "排序重复"
		return nil, errors.New(msg)
	}
	elementID, err := repo.CreateElement(info)
	if err != nil {
		return nil, err
	}
	element, err := repo.GetElementByID(elementID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return element, err
}

func (s *elementService) GetElementList(filter ElementFilter) (int, *[]Element, error) {
	db := database.InitMySQL()
	query := NewElementQuery(db)
	count, err := query.GetElementCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetElementList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *elementService) UpdateElement(elementID int64, info ElementUpdate, organizationID int64) (*Element, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewElementRepository(tx)
	oldElement, err := repo.GetElementByID(elementID)
	if err != nil {
		return nil, err
	}
	nodeExist, err := repo.CheckNodeExist(oldElement.NodeID, organizationID)
	if err != nil {
		return nil, err
	}
	if nodeExist == 0 {
		msg := "节点不存在"
		return nil, errors.New(msg)
	}
	if info.Name != "" {
		exist, err := repo.CheckNameExist(info.Name, oldElement.NodeID, elementID)
		if err != nil {
			return nil, err
		}
		if exist != 0 {
			msg := "节点名称重复"
			return nil, errors.New(msg)
		}
		oldElement.Name = info.Name
	}
	if info.Type != "" {
		oldElement.ElementType = info.Type
	}
	if info.Sort != 0 {
		exist, err := repo.CheckSortExist(info.Sort, oldElement.NodeID, elementID)
		if err != nil {
			return nil, err
		}
		if exist != 0 {
			msg := "排序重复"
			return nil, errors.New(msg)
		}
		oldElement.Sort = info.Sort
	}
	if info.Required != 0 {
		oldElement.Required = info.Required
	}
	if info.Patterns != "" {
		oldElement.Patterns = info.Patterns
	}
	if info.DefaultValue != "" {
		oldElement.DefaultValue = info.DefaultValue
	}
	oldElement.JsonData = info.JsonData

	err = repo.UpdateElement(elementID, *oldElement, info.User)
	if err != nil {
		return nil, err
	}
	element, err := repo.GetElementByID(elementID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return element, err
}

func (s *elementService) DeleteElement(elementID int64, user string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewElementRepository(tx)
	err = repo.DeleteElement(elementID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
