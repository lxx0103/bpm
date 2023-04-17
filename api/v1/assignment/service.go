package assignment

import (
	"bpm/api/v1/auth"
	"bpm/api/v1/member"
	"bpm/api/v1/project"
	"bpm/core/database"
	"errors"
	"fmt"
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
	projectRepo := project.NewProjectRepository(tx)
	memberRepo := member.NewMemberRepository(tx)
	userRepo := auth.NewAuthRepository(tx)
	_, err = projectRepo.GetProjectByID(info.ProjectID, info.OrganizationID)
	if err != nil {
		msg := "获取项目失败"
		return errors.New(msg)
	}
	memberExist, err := memberRepo.CheckMemberExist(info.ProjectID, info.AssignTo)
	if err != nil {
		fmt.Println(info.ProjectID, info.AssignTo, err.Error())
		msg := "获取项目成员失败"
		return errors.New(msg)
	}
	if !memberExist {
		msg := "只能把任务分配给项目成员"
		return errors.New(msg)
	}
	user, err := userRepo.GetUserByID(info.AuditTo)
	if err != nil {
		msg := "获取审核人员失败"
		return errors.New(msg)
	}
	if user.OrganizationID != info.OrganizationID {
		msg := "审核人员不存在"
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

func (s *assignmentService) UpdateAssignment(assignmentID int64, info AssignmentUpdate, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewAssignmentRepository(tx)
	projectRepo := project.NewProjectRepository(tx)
	memberRepo := member.NewMemberRepository(tx)
	userRepo := auth.NewAuthRepository(tx)
	oldAssignment, err := repo.GetAssignmentByID(assignmentID)
	if err != nil {
		msg := "任务记录不存在"
		return errors.New(msg)
	}
	if organizationID != oldAssignment.OrganizationID && organizationID != 0 {
		msg := "任务记录不存在"
		return errors.New(msg)
	}
	if oldAssignment.Status == 9 {
		msg := "此任务已完成"
		return errors.New(msg)
	}
	if oldAssignment.CreatedBy != info.User {
		msg := "只能修改自己创建的任务"
		return errors.New(msg)
	}
	_, err = projectRepo.GetProjectByID(info.ProjectID, oldAssignment.OrganizationID)
	if err != nil {
		msg := "获取项目失败"
		return errors.New(msg)
	}
	memberExist, err := memberRepo.CheckMemberExist(info.ProjectID, info.AssignTo)
	if err != nil {
		fmt.Println(info.ProjectID, info.AssignTo, err.Error())
		msg := "获取项目成员失败"
		return errors.New(msg)
	}
	if !memberExist {
		msg := "只能把任务分配给项目成员"
		return errors.New(msg)
	}
	user, err := userRepo.GetUserByID(info.AuditTo)
	if err != nil {
		msg := "获取审核人员失败"
		return errors.New(msg)
	}
	if user.OrganizationID != oldAssignment.OrganizationID {
		msg := "审核人员不存在"
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
	if oldAssignment.CreatedBy != byUser {
		msg := "只能删除自己创建的任务"
		return errors.New(msg)
	}
	err = repo.DeleteAssignment(assignmentID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *assignmentService) CompleteAssignment(assignmentID int64, info AssignmentComplete, organizationID int64) error {
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
	if oldAssignment.Status != 1 && oldAssignment.Status != 3 {
		msg := "此任务不可完成"
		return errors.New(msg)
	}
	if oldAssignment.AssignTo != info.UserID {
		msg := "只能完成分配给你的任务"
		return errors.New(msg)
	}
	err = repo.CompleteAssignment(assignmentID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *assignmentService) AuditAssignment(assignmentID int64, info AssignmentAudit, organizationID int64) error {
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
	if oldAssignment.Status != 2 {
		msg := "此任务不可审核"
		return errors.New(msg)
	}
	if oldAssignment.AuditTo != info.UserID {
		msg := "只能审核分配给你的任务"
		return errors.New(msg)
	}
	err = repo.AuditAssignment(assignmentID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *assignmentService) GetMyAssignmentList(filter MyAssignmentFilter) (int, *[]AssignmentResponse, error) {
	db := database.InitMySQL()
	query := NewAssignmentQuery(db)
	count, err := query.GetMyAssignmentCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetMyAssignmentList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *assignmentService) GetMyAuditList(filter MyAuditFilter) (int, *[]AssignmentResponse, error) {
	db := database.InitMySQL()
	query := NewAssignmentQuery(db)
	count, err := query.GetMyAuditCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetMyAuditList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}
