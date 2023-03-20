package assignment

import (
	"bpm/core/database"
	"errors"
)

type assignmentService struct {
}

func NewAssignmentService() *assignmentService {
	return &assignmentService{}
}

func (s *assignmentService) GetAssignmentByID(id int64, organizationID int64) (*AssignmentResponse, error) {
	db := database.InitMySQL()
	query := NewAssignmentQuery(db)
	assignment, err := query.GetAssignmentByID(id, organizationID)
	return assignment, err
}

func (s *assignmentService) NewAssignment(info AssignmentNew, organizationID int64) error {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewAssignmentRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, 0)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "任务名称重复"
		return errors.New(msg)
	}
	err = repo.CreateAssignment(info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *assignmentService) GetAssignmentList(filter AssignmentFilter, organizationID int64) (int, *[]AssignmentResponse, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewAssignmentQuery(db)
	count, err := query.GetAssignmentCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetAssignmentList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *assignmentService) UpdateAssignment(assignmentID int64, info AssignmentNew, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewAssignmentRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, assignmentID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "任务记录名称重复"
		return errors.New(msg)
	}
	oldAssignment, err := repo.GetAssignmentByID(assignmentID)
	if err != nil {
		msg := "任务记录不存在"
		return errors.New(msg)
	}
	if organizationID != oldAssignment.OrganizationID && organizationID != 0 {
		msg := "任务记录不存在"
		return errors.New(msg)
	}
	err = repo.UpdateAssignment(assignmentID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *assignmentService) DeleteAssignment(assignmentID, organizationID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewAssignmentRepository(tx)
	oldAssignment, err := repo.GetAssignmentByID(assignmentID)
	if err != nil {
		msg := "任务记录不存在"
		return errors.New(msg)
	}
	if organizationID != oldAssignment.OrganizationID && organizationID != 0 {
		msg := "任务记录不存在"
		return errors.New(msg)
	}
	err = repo.DeleteAssignment(assignmentID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
