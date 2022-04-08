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
	GetComponentList(ComponentFilter) (int, *[]Component, error)
}

func (s *componentService) GetComponentByID(id int64) (*Component, error) {
	db := database.InitMySQL()
	query := NewComponentQuery(db)
	component, err := query.GetComponentByID(id)
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
