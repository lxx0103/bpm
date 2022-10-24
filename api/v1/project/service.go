package project

import (
	"bpm/api/v1/component"
	"bpm/api/v1/element"
	"bpm/api/v1/event"
	"bpm/api/v1/member"
	"bpm/api/v1/node"
	"bpm/api/v1/template"
	"bpm/core/database"
	"bpm/core/queue"
	"encoding/json"
	"errors"
	"fmt"
)

type projectService struct {
}

func NewProjectService() *projectService {
	return &projectService{}
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
	memberRepo := member.NewMemberRepository(tx)
	template, err := templateRepo.GetTemplateByID(info.TemplateID)
	var projectMember []int64
	// projectMember = append(projectMember, info.UserID)
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
	info.Type = template.Type
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
		eventInfo.NeedAudit = (*nodes)[i].NeedAudit
		eventInfo.AuditType = (*nodes)[i].AuditType
		eventInfo.NeedCheckin = (*nodes)[i].NeedCheckin
		eventInfo.Sort = (*nodes)[i].Sort
		eventInfo.CanReview = (*nodes)[i].CanReview
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
	events, err := eventRepo.GetEventsByProjectID(projectID)
	if err != nil {
		return nil, err
	}
	for k := 0; k < len(*events); k++ {
		var pres []int64
		var assigns []int64
		var audits []int64
		nodePres, err := nodeRepo.GetPresByNodeID((*events)[k].NodeID)
		if err != nil {
			return nil, err
		}
		for l := 0; l < len(*nodePres); l++ {
			preEventID, err := eventRepo.GetEventIDByProjectAndNode(projectID, (*nodePres)[l].PreID)
			if err != nil {
				return nil, err
			}
			pres = append(pres, preEventID)
		}
		err = eventRepo.CreateEventPre((*events)[k].ID, pres, info.User)
		if err != nil {
			return nil, err
		}
		nodeAudits, err := nodeRepo.GetAuditsByNodeID((*events)[k].NodeID)
		if err != nil {
			return nil, err
		}
		for n := 0; n < len(*nodeAudits); n++ {
			audits = append(audits, (*nodeAudits)[n].AuditTo)
			if (*nodeAudits)[n].AuditType == 2 {
				projectMember = append(projectMember, (*nodeAudits)[n].AuditTo)
			}
		}
		err = eventRepo.CreateEventAudit((*events)[k].ID, (*events)[k].AuditType, audits, info.User)
		if err != nil {
			return nil, err
		}
		if (*events)[k].AssignType == 3 {
			assigns = append(assigns, info.UserID)
			projectMember = append(projectMember, info.UserID)
			(*events)[k].AssignType = 2
		} else {
			nodeAssigns, err := nodeRepo.GetAssignsByNodeID((*events)[k].NodeID)
			if err != nil {
				return nil, err
			}
			for m := 0; m < len(*nodeAssigns); m++ {
				assigns = append(assigns, (*nodeAssigns)[m].AssignTo)
				if (*nodeAssigns)[m].AssignType == 2 {
					projectMember = append(projectMember, (*nodeAssigns)[m].AssignTo)
				}
			}
		}
		err = eventRepo.CreateEventAssign((*events)[k].ID, (*events)[k].AssignType, assigns, info.User)
		if err != nil {
			return nil, err
		}
		err = eventRepo.UpdateEvent((*events)[k].ID, (*events)[k], info.User)
		if err != nil {
			return nil, err
		}
	}
	err = memberRepo.DeleteProjectMember(projectID, info.User)
	if err != nil {
		return nil, err
	}
	memberRepo.CreateProjectMember(projectID, projectMember, organizationID, info.User)
	project, err := repo.GetProjectByID(projectID, organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	type NewProjectCreated struct {
		ProjectID int64 `json:"project_id"`
	}
	var newEvent NewProjectCreated
	newEvent.ProjectID = projectID
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewProjectCreated", msg)
	if err != nil {
		msg := "create event NewProjectCreated error"
		return nil, errors.New(msg)
	}
	return project, err
}

func (s *projectService) GetProjectList(filter ProjectFilter, organizationID int64) (int, *[]ProjectResponse, error) {
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

func (s *projectService) UpdateProject(projectID int64, info ProjectUpdate, organizationID int64) (*Project, error) {
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
	if info.Name != "" {
		exist, err := repo.CheckNameExist(info.Name, organizationID, projectID)
		if err != nil {
			return nil, err
		}
		if exist != 0 {
			msg := "项目名称重复"
			return nil, errors.New(msg)
		}
		oldProject.Name = info.Name
	}
	if info.ClientID != 0 {
		oldProject.ClientID = info.ClientID
	}
	if info.Location != "" {
		oldProject.Location = info.Location
	}
	if info.Longitude != 0 {
		oldProject.Longitude = info.Longitude
	}
	if info.Latitude != 0 {
		oldProject.Latitude = info.Latitude
	}
	if info.CheckinDistance != 0 {
		oldProject.CheckinDistance = info.CheckinDistance
	}
	if info.Priority != 0 {
		oldProject.Priority = info.Priority
	}
	err = repo.UpdateProject(projectID, *oldProject, info.User)
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

func (s *projectService) DeleteProject(projectID int64, organizationID int64, user string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	eventRepo := event.NewEventRepository(tx)
	oldProject, err := repo.GetProjectByID(projectID, organizationID)
	if err != nil {
		return err
	}
	if oldProject.CreatedBy != user {
		msg := "只能删除你创建的项目"
		return errors.New(msg)
	}
	err = repo.DeleteProject(projectID, user)
	if err != nil {
		return err
	}
	err = eventRepo.DeleteEventByProjectID(projectID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *projectService) GetMyProject(filter MyProjectFilter, userName string, organizationID int64) (int, *[]Project, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)

	myProjects, err := query.GetProjectListByCreate(userName, organizationID, filter)
	if err != nil {
		return 0, nil, err
	}
	myProjectsCount, err := query.GetProjectCountByCreate(userName, organizationID, filter)
	if err != nil {
		return 0, nil, err
	}
	return myProjectsCount, myProjects, err
}

func (s *projectService) GetAssignedProject(filter AssignedProjectFilter, userID int64, positionID int64, organizationID int64) (int, *[]Project, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)

	myProjects, err := query.GetProjectListByAssigned(filter, userID, positionID, organizationID)
	if err != nil {
		fmt.Println("aaa")
		return 0, nil, err
	}
	myProjectsCount, err := query.GetProjectCountByAssigned(filter, userID, positionID, organizationID)
	if err != nil {
		fmt.Println("bbb")
		return 0, nil, err
	}
	return myProjectsCount, myProjects, err
}

func (s *projectService) GetClientProject(filter MyProjectFilter, userID int64, organizationID int64) (int, *[]Project, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)

	myProjects, err := query.GetProjectListByClientID(userID, organizationID, filter)
	if err != nil {
		return 0, nil, err
	}
	myProjectsCount, err := query.GetProjectCountByClientID(userID, organizationID, filter)
	if err != nil {
		return 0, nil, err
	}
	return myProjectsCount, myProjects, err
}
