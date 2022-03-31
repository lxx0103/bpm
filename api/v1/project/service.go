package project

import (
	"bpm/api/v1/component"
	"bpm/api/v1/element"
	"bpm/api/v1/event"
	"bpm/api/v1/node"
	"bpm/api/v1/template"
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
	templateRepo := template.NewTemplateRepository(tx)
	nodeRepo := node.NewNodeRepository(tx)
	elementRepo := element.NewElementRepository(tx)
	eventRepo := event.NewEventRepository(tx)
	componentRepo := component.NewComponentRepository(tx)
	template, err := templateRepo.GetTemplateByID(info.TemplateID)
	if err != nil {
		return nil, err
	}
	if organizationID != 0 && template.OrganizationID != organizationID {
		msg := "你无权使用此模板"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, organizationID, 0)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "项目名称重复"
		return nil, errors.New(msg)
	}
	projectID, err := repo.CreateProject(info, template.OrganizationID)
	if err != nil {
		return nil, err
	}
	nodes, err := nodeRepo.GetNodesByTemplateID(info.TemplateID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(*nodes); i++ {
		var eventInfo event.EventNew
		eventInfo.ProjectID = projectID
		eventInfo.Name = (*nodes)[i].Name
		eventInfo.AssignType = (*nodes)[i].AssignType
		eventInfo.Assignable = (*nodes)[i].Assignable
		eventInfo.NodeID = (*nodes)[i].ID
		eventInfo.User = info.User
		eventID, err := eventRepo.CreateEvent(eventInfo)
		if err != nil {
			return nil, err
		}
		elements, err := elementRepo.GetElementsByNodeID((*nodes)[i].ID)
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(*elements); j++ {
			var componentInfo component.ComponentNew
			componentInfo.EventID = eventID
			componentInfo.Sort = (*elements)[j].Sort
			componentInfo.Type = (*elements)[j].ElementType
			componentInfo.Name = (*elements)[j].Name
			componentInfo.DefaultValue = (*elements)[j].DefaultValue
			componentInfo.Required = (*elements)[j].Required
			componentInfo.Patterns = (*elements)[j].Patterns
			componentInfo.User = info.User
			_, err := componentRepo.CreateComponent(componentInfo)
			if err != nil {
				return nil, err
			}
		}
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
		msg := "你无权修改此项目"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, organizationID, projectID)
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
