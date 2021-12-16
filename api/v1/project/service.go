package project

import (
	"bpm/core/database"
)

type projectService struct {
}

func NewProjectService() ProjectService {
	return &projectService{}
}

// ProjectService represents a service for managing projects.
type ProjectService interface {
	//Project Management
	GetProjectByID(int64) (*Project, error)
	NewProject(ProjectNew) (*Project, error)
	GetProjectList(ProjectFilter) (int, *[]Project, error)
	UpdateProject(int64, ProjectNew) (*Project, error)
}

func (s *projectService) GetProjectByID(id int64) (*Project, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	project, err := query.GetProjectByID(id)
	return project, err
}

func (s *projectService) NewProject(info ProjectNew) (*Project, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	projectID, err := repo.CreateProject(info)
	if err != nil {
		return nil, err
	}
	project, err := repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return project, err
}

func (s *projectService) GetProjectList(filter ProjectFilter) (int, *[]Project, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	count, err := query.GetProjectCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetProjectList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *projectService) UpdateProject(projectID int64, info ProjectNew) (*Project, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	_, err = repo.UpdateProject(projectID, info)
	if err != nil {
		return nil, err
	}
	project, err := repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return project, err
}
