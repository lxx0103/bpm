package component

import (
	"bpm/core/database"
)

type componentService struct {
}

func NewComponentService() ComponentService {
	return &componentService{}
}

// ComponentService represents a service for managing components.
type ComponentService interface {
	//Component Management
	GetComponentByID(int64) (*Component, error)
	NewComponent(ComponentNew) (*Component, error)
	GetComponentList(ComponentFilter) (int, *[]Component, error)
	UpdateComponent(int64, ComponentUpdate) (*Component, error)
	DeleteComponent(int64, string) error
}

func (s *componentService) GetComponentByID(id int64) (*Component, error) {
	db := database.InitMySQL()
	query := NewComponentQuery(db)
	component, err := query.GetComponentByID(id)
	return component, err
}

func (s *componentService) NewComponent(info ComponentNew) (*Component, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewComponentRepository(tx)
	componentID, err := repo.CreateComponent(info)
	if err != nil {
		return nil, err
	}
	component, err := repo.GetComponentByID(componentID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return component, err
}

func (s *componentService) GetComponentList(filter ComponentFilter) (int, *[]Component, error) {
	db := database.InitMySQL()
	query := NewComponentQuery(db)
	count, err := query.GetComponentCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetComponentList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *componentService) UpdateComponent(componentID int64, info ComponentUpdate) (*Component, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewComponentRepository(tx)
	oldComponent, err := repo.GetComponentByID(componentID)
	if err != nil {
		return nil, err
	}
	if info.Type != "" {
		oldComponent.ComponentType = info.Type
	}
	if info.Sort != 0 {
		oldComponent.Sort = info.Sort
	}
	if info.Name != "" {
		oldComponent.Name = info.Name
	}
	if info.Required != 0 {
		oldComponent.Required = info.Required
	}
	if info.Patterns != "" {
		oldComponent.Patterns = info.Patterns
	}
	if info.DefaultValue != "" {
		oldComponent.DefaultValue = info.DefaultValue
	}
	oldComponent.JsonData = info.JsonData

	err = repo.UpdateComponent(componentID, *oldComponent, info.User)
	if err != nil {
		return nil, err
	}
	component, err := repo.GetComponentByID(componentID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return component, err
}

func (s *componentService) DeleteComponent(componentID int64, user string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewComponentRepository(tx)
	err = repo.DeleteComponent(componentID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
