package position

import (
	"bpm/core/database"
	"errors"
)

type positionService struct {
}

func NewPositionService() PositionService {
	return &positionService{}
}

// PositionService represents a service for managing positions.
type PositionService interface {
	//Position Management
	GetPositionByID(int64, int64) (*Position, error)
	NewPosition(PositionNew, int64) (*Position, error)
	GetPositionList(PositionFilter, int64) (int, *[]PositionResponse, error)
	UpdatePosition(int64, PositionNew, int64) (*Position, error)
}

func (s *positionService) GetPositionByID(id int64, organizationID int64) (*Position, error) {
	db := database.InitMySQL()
	query := NewPositionQuery(db)
	position, err := query.GetPositionByID(id, organizationID)
	return position, err
}

func (s *positionService) NewPosition(info PositionNew, organizationID int64) (*Position, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewPositionRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, organizationID, 0)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "职位名称重复"
		return nil, errors.New(msg)
	}
	positionID, err := repo.CreatePosition(info, organizationID)
	if err != nil {
		return nil, err
	}
	position, err := repo.GetPositionByID(positionID, organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return position, err
}

func (s *positionService) GetPositionList(filter PositionFilter, organizationID int64) (int, *[]PositionResponse, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewPositionQuery(db)
	count, err := query.GetPositionCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPositionList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *positionService) UpdatePosition(positionID int64, info PositionNew, organizationID int64) (*Position, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewPositionRepository(tx)
	oldPosition, err := repo.GetPositionByID(positionID, organizationID)
	if err != nil {
		return nil, err
	}
	if organizationID != 0 && organizationID != oldPosition.OrganizationID {
		msg := "你无权修改此职位"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, organizationID, positionID)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "职位名称重复"
		return nil, errors.New(msg)
	}
	_, err = repo.UpdatePosition(positionID, info)
	if err != nil {
		return nil, err
	}
	position, err := repo.GetPositionByID(positionID, organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return position, err
}
