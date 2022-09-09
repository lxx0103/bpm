package example

import (
	"bpm/core/database"
	"errors"
)

type exampleService struct {
}

func NewExampleService() ExampleService {
	return &exampleService{}
}

// ExampleService represents a service for managing examples.
type ExampleService interface {
	//Example Management
	GetExampleByID(int64, int64) (*Example, error)
	NewExample(ExampleNew, int64) (*Example, error)
	GetExampleList(ExampleFilter, int64) (int, *[]ExampleListResponse, error)
	UpdateExample(int64, ExampleNew, int64) (*Example, error)
}

func (s *exampleService) GetExampleByID(id int64, organizationID int64) (*Example, error) {
	db := database.InitMySQL()
	query := NewExampleQuery(db)
	example, err := query.GetExampleByID(id, organizationID)
	return example, err
}

func (s *exampleService) NewExample(info ExampleNew, organizationID int64) (*Example, error) {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewExampleRepository(tx)
	// exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, 0)
	// if err != nil {
	// 	return nil, err
	// }
	// if exist != 0 {
	// 	msg := "案例名称重复"
	// 	return nil, errors.New(msg)
	// }
	exampleID, err := repo.CreateExample(info)
	if err != nil {
		return nil, err
	}
	example, err := repo.GetExampleByID(exampleID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return example, err
}

func (s *exampleService) GetExampleList(filter ExampleFilter, organizationID int64) (int, *[]ExampleListResponse, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewExampleQuery(db)
	count, err := query.GetExampleCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetExampleList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *exampleService) UpdateExample(exampleID int64, info ExampleNew, organizationID int64) (*Example, error) {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewExampleRepository(tx)
	// exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, exampleID)
	// if err != nil {
	// 	return nil, err
	// }
	// if exist != 0 {
	// 	msg := "案例名称重复"
	// 	return nil, errors.New(msg)
	// }
	_, err = repo.UpdateExample(exampleID, info)
	if err != nil {
		return nil, err
	}
	example, err := repo.GetExampleByID(exampleID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return example, err
}
