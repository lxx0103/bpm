package project

import (
	"bpm/core/database"
	"errors"
)

type projectService struct {
}

func NewProjectService() ProjectService {
	return &projectService{}
}

// ProjectService represents a service for managing projects.
type ProjectService interface {
	//Project Management
	GetProjectByID(int64, int64) (*Project, error)
	NewProject(ProjectNew, int64) (*Project, error)
	GetProjectList(ProjectFilter, int64) (int, *[]Project, error)
	UpdateProject(int64, ProjectNew, int64) (*Project, error)
}

func (s *projectService) GetProjectByID(id int64, organizationID int64) (*Project, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	project, err := query.GetProjectByID(id, organizationID)
	return project, err
}

func (s *projectService) NewProject(info ProjectNew, organizationID int64) (*Project, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, organizationID)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "项目名称重复"
		return nil, errors.New(msg)
	}
	projectID, err := repo.CreateProject(info, organizationID)
	if err != nil {
		return nil, err
	}
	project, err := repo.GetProjectByID(projectID, organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return project, err
}

func (s *projectService) GetProjectList(filter ProjectFilter, organizationID int64) (int, *[]Project, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	count, err := query.GetProjectCount(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetProjectList(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *projectService) UpdateProject(projectID int64, info ProjectNew, organizationID int64) (*Project, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	oldProject, err := repo.GetProjectByID(projectID, organizationID)
	if err != nil {
		return nil, err
	}
	if organizationID != 0 && organizationID != oldProject.OrganizationID {
		msg := "你无权修改此客户"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, organizationID)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "项目名称重复"
		return nil, errors.New(msg)
	}
	_, err = repo.UpdateProject(projectID, info)
	if err != nil {
		return nil, err
	}
	project, err := repo.GetProjectByID(projectID, organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return project, err
}
