package project

import (
	"bpm/api/v1/component"
	"bpm/api/v1/element"
	"bpm/api/v1/event"
	"bpm/api/v1/member"
	"bpm/api/v1/node"
	"bpm/api/v1/team"
	"bpm/api/v1/template"
	"bpm/core/database"
	"bpm/core/queue"
	"encoding/json"
	"errors"
	"fmt"
	"time"
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
	if err != nil {
		msg := "获取项目失败"
		return nil, errors.New(msg)
	}
	teams, err := query.GetProjectTeam(id)
	if err != nil {
		msg := "获取项目班组失败"
		return nil, errors.New(msg)
	}
	project.Teams = *teams
	return project, nil
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
		if len(*nodePres) == 0 {
			err = eventRepo.SetEventActive((*events)[k].ID)
			if err != nil {
				return nil, err
			}
		}
		nodeAudits, err := nodeRepo.GetAuditsByNodeID((*events)[k].NodeID)
		if err != nil {
			return nil, err
		}
		for n := 0; n < len(*nodeAudits); n++ {
			var nodeAudit event.NodeAudit
			nodeAudit.AuditLevel = (*nodeAudits)[n].AuditLevel
			nodeAudit.AuditTo = append(nodeAudit.AuditTo, (*nodeAudits)[n].AuditTo)
			err = eventRepo.CreateEventAudit((*events)[k].ID, (*nodeAudits)[n].AuditType, nodeAudit, info.User)
			if err != nil {
				return nil, err
			}
			if (*nodeAudits)[n].AuditType == 2 {
				projectMember = append(projectMember, (*nodeAudits)[n].AuditTo)
			}
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
	if len(info.TeamID) > 0 {
		teamRepo := team.NewTeamRepository(tx)
		for _, teamID := range info.TeamID {
			_, err = teamRepo.GetTeamByID(teamID, organizationID)
			if err != nil {
				msg := "班组不存在"
				return nil, errors.New(msg)
			}
			err = repo.CreateProjectTeam(projectID, teamID, info.User)
			if err != nil {
				msg := "创建班组信息失败"
				return nil, errors.New(msg)
			}
		}
	}
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
	for k, v := range *list {
		teams, err := query.GetProjectTeam(v.ID)
		if err != nil {
			msg := "获取项目班组失败" + err.Error()
			return 0, nil, errors.New(msg)
		}
		(*list)[k].Teams = *teams
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
	if info.Area != "" {
		oldProject.Area = info.Area
	}
	oldProject.RecordAlertDay = info.RecordAlertDay
	err = repo.UpdateProject(projectID, *oldProject, info.User)
	if err != nil {
		return nil, err
	}
	err = repo.DeleteProjectTeam(projectID, info.User)
	if err != nil {
		return nil, err
	}
	if len(info.TeamID) > 0 {
		teamRepo := team.NewTeamRepository(tx)
		for _, teamID := range info.TeamID {
			_, err = teamRepo.GetTeamByID(teamID, organizationID)
			if err != nil {
				msg := "班组不存在"
				return nil, errors.New(msg)
			}
			err = repo.CreateProjectTeam(projectID, teamID, info.User)
			if err != nil {
				msg := "创建班组信息失败"
				return nil, errors.New(msg)
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
	err = repo.DeleteReportByProjectID(projectID, user)
	if err != nil {
		return err
	}
	err = repo.DeleteRecordByProjectID(projectID, user)
	if err != nil {
		return err
	}
	err = repo.DeleteAssignmentByProjectID(projectID, user)
	if err != nil {
		return err
	}
	err = repo.DeleteProjectTeam(projectID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *projectService) GetMyProject(filter MyProjectFilter, userName string, organizationID int64) (int, *[]ProjectResponse, error) {
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
	for k, v := range *myProjects {
		events, err := query.GetActiveEvents(v.ID)
		if err != nil {
			return 0, nil, err
		}
		for k2, v2 := range *events {
			if v2.Status == 1 || v2.Status == 3 {
				(*events)[k2].ActiveType = "执行"
				if v2.AssignType == 1 {
					assigns, err := query.GetEventAssignPosition(v2.EventID)
					if err != nil {
						return 0, nil, err
					}
					(*events)[k2].Actives = *assigns
				} else {
					assigns, err := query.GetEventAssignUser(v2.EventID)
					if err != nil {
						return 0, nil, err
					}
					(*events)[k2].Actives = *assigns
				}
			} else if v2.Status == 2 {
				(*events)[k2].ActiveType = "审核"
				// if v2.AuditType == 1 {
				assigns, err := query.GetEventAuditPosition(v2.EventID)
				if err != nil {
					return 0, nil, err
				}
				(*events)[k2].Actives = *assigns
				// } else {
				// 	assigns, err := query.GetEventAuditUser(v2.EventID)
				// 	if err != nil {
				// 		return 0, nil, err
				// 	}
				// 	(*events)[k2].Actives = *assigns
				// }
			}
		}
		(*myProjects)[k].ActiveEvents = *events
	}
	return myProjectsCount, myProjects, err
}

func (s *projectService) GetAssignedProject(filter AssignedProjectFilter, userID int64, positionID int64, organizationID int64) (int, *[]ProjectResponse, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)

	myProjects, err := query.GetProjectListByAssigned(filter, userID, positionID, organizationID)
	if err != nil {
		return 0, nil, err
	}
	myProjectsCount, err := query.GetProjectCountByAssigned(filter, userID, positionID, organizationID)
	if err != nil {
		return 0, nil, err
	}
	for k, v := range *myProjects {
		events, err := query.GetActiveEvents(v.ID)
		if err != nil {
			return 0, nil, err
		}
		for k2, v2 := range *events {
			if v2.Status == 1 || v2.Status == 3 {
				(*events)[k2].ActiveType = "执行"
				if v2.AssignType == 1 {
					assigns, err := query.GetEventAssignPosition(v2.EventID)
					if err != nil {
						return 0, nil, err
					}
					(*events)[k2].Actives = *assigns
				} else {
					assigns, err := query.GetEventAssignUser(v2.EventID)
					if err != nil {
						return 0, nil, err
					}
					(*events)[k2].Actives = *assigns
				}
			} else if v2.Status == 2 {
				(*events)[k2].ActiveType = "审核"
				// if v2.AuditType == 1 {
				assigns, err := query.GetEventAuditPosition(v2.EventID)
				if err != nil {
					return 0, nil, err
				}
				(*events)[k2].Actives = *assigns
				// } else {
				// 	assigns, err := query.GetEventAuditUser(v2.EventID)
				// 	if err != nil {
				// 		return 0, nil, err
				// 	}
				// 	(*events)[k2].Actives = *assigns
				// }
			}
		}
		(*myProjects)[k].ActiveEvents = *events

		startDate := v.Created.Format("2006-01-02")
		lastRecordDate := v.LastRecordDate
		if v.LastRecordDate == "" {
			lastRecordDate = startDate
		}
		recordDateTime, err := time.Parse("2006-01-02", lastRecordDate)
		if err != nil {
			return 0, nil, err
		}
		today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		if err != nil {
			return 0, nil, err
		}
		interval := today.Sub(recordDateTime)
		(*myProjects)[k].NoRecordDay = int(interval.Hours() / 24)
	}
	return myProjectsCount, myProjects, err
}

func (s *projectService) GetClientProject(filter MyProjectFilter, userID int64, organizationID int64) (int, *[]ProjectResponse, error) {
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
	for k, v := range *myProjects {
		events, err := query.GetActiveEvents(v.ID)
		if err != nil {
			return 0, nil, err
		}
		for k2, v2 := range *events {
			if v2.Status == 1 || v2.Status == 3 {
				(*events)[k2].ActiveType = "执行"
				if v2.AssignType == 1 {
					assigns, err := query.GetEventAssignPosition(v2.EventID)
					if err != nil {
						return 0, nil, err
					}
					(*events)[k2].Actives = *assigns
				} else {
					assigns, err := query.GetEventAssignUser(v2.EventID)
					if err != nil {
						return 0, nil, err
					}
					(*events)[k2].Actives = *assigns
				}
			} else if v2.Status == 2 {
				(*events)[k2].ActiveType = "审核"
				// if v2.AuditType == 1 {
				assigns, err := query.GetEventAuditPosition(v2.EventID)
				if err != nil {
					return 0, nil, err
				}
				(*events)[k2].Actives = *assigns
				// } else {
				// 	assigns, err := query.GetEventAuditUser(v2.EventID)
				// 	if err != nil {
				// 		return 0, nil, err
				// 	}
				// 	(*events)[k2].Actives = *assigns
				// }
			}
		}
		(*myProjects)[k].ActiveEvents = *events
	}
	return myProjectsCount, myProjects, err
}

func (s *projectService) NewProjectReport(projectID int64, info ProjectReportNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	memberRepo := member.NewMemberRepository(tx)
	project, err := repo.GetProjectByID(projectID, info.OrganizationID)
	if err != nil {
		msg := "项目不存在"
		return errors.New(msg)
	}
	members, err := memberRepo.GetMembersByProjectID(projectID)
	if err != nil {
		msg := "获取项目成员失败"
		return errors.New(msg)
	}
	memberValid := false
	for _, member := range *members {
		if member.UserID == info.UserID {
			memberValid = true
			break
		}
	}
	if !memberValid {
		msg := "你不是此项目的成员"
		return errors.New(msg)
	}
	var newReport ProjectReport
	newReport.OrganizationID = info.OrganizationID
	newReport.ProjectID = projectID
	newReport.ClientID = project.ClientID
	newReport.UserID = info.UserID
	newReport.Content = info.Content
	newReport.ReportDate = info.ReportDate
	newReport.Name = info.Name
	newReport.Content = info.Content
	newReport.Status = 1
	newReport.Created = time.Now()
	newReport.CreatedBy = info.User
	newReport.Updated = time.Now()
	newReport.UpdatedBy = info.User
	reportID, err := repo.CreateProjectReport(newReport)
	if err != nil {
		msg := "创建报告失败"
		return errors.New(msg)
	}
	for _, link := range info.Links {
		var reportLink ProjectReportLink
		reportLink.OrganizationID = info.OrganizationID
		reportLink.ProjectID = projectID
		reportLink.ProjectReportID = reportID
		reportLink.Link = link
		reportLink.Status = 1
		reportLink.Created = time.Now()
		reportLink.CreatedBy = info.User
		reportLink.Updated = time.Now()
		reportLink.UpdatedBy = info.User
		err = repo.CreateProjectReportLink(reportLink)
		if err != nil {
			msg := "创建链接失败"
			return errors.New(msg)
		}
	}

	var newReportView ProjectReportView
	newReportView.OrganizationID = info.OrganizationID
	newReportView.ProjectID = projectID
	newReportView.ProjectReportID = reportID
	newReportView.ViewerID = info.UserID
	newReportView.ViewerName = info.User
	newReportView.Status = 1
	newReportView.Created = time.Now()
	newReportView.CreatedBy = info.User
	newReportView.Updated = time.Now()
	newReportView.UpdatedBy = info.User
	err = repo.CreateProjectReportView(newReportView)
	if err != nil {
		msg := "创建已阅失败"
		return errors.New(msg)
	}
	views, err := repo.GetProjectReportView(reportID)
	if err != nil {
		msg := "获取阅读记录失败"
		return errors.New(msg)
	}
	reportStatus := 1
	if len(*views) == len(*members) {
		reportStatus = 3
	} else {
		reportStatus = 2
	}
	err = repo.UpdateProjectReportStatus(reportID, reportStatus, info.User)
	if err != nil {
		msg := "更新状态失败"
		return errors.New(msg)
	}
	tx.Commit()

	type NewProjectReportCreated struct {
		ProjectReportID int64 `json:"project_report_id"`
	}
	var newEvent NewProjectReportCreated
	newEvent.ProjectReportID = reportID
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewProjectReportCreated", msg)
	if err != nil {
		msg := "发布消息（NewProjectReportCreated）失败"
		return errors.New(msg)
	}
	return err
}

func (s *projectService) GetProjectReportList(projectID int64, filter ProjectReportFilter) (*[]ProjectReportResponse, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	memberQuery := member.NewMemberQuery(db)
	_, err := query.GetProjectByID(projectID, filter.OrganizationID)
	if err != nil {
		msg := "项目不存在"
		return nil, errors.New(msg)
	}
	members, err := memberQuery.GetMembersByProjectID(projectID)
	if err != nil {
		msg := "获取成员失败" + err.Error()
		return nil, errors.New(msg)
	}
	memberValid := false
	for _, member := range *members {
		if member.UserID == filter.UserID {
			memberValid = true
			break
		}
	}
	if filter.OrganizationID == 0 {
		memberValid = true
	}
	if !memberValid {
		msg := "你不是此项目的成员"
		return nil, errors.New(msg)
	}
	list, err := query.GetProjectReportList(projectID, filter)

	for k, v := range *list {
		(*list)[k].Viewed = false
		links, err := query.GetProjectReportLinks(v.ID)
		if err != nil {
			msg := "获取报告链接失败" // + err.Error()
			return nil, errors.New(msg)
		}
		(*list)[k].Links = *links
		views, err := query.GetProjectReportViews(v.ID)
		if err != nil {
			msg := "获取报告阅读记录失败" // + err.Error()
			return nil, errors.New(msg)
		}
		var memberViews []ProjectReportMemberViewResponse
		for _, member := range *members {
			var memberView ProjectReportMemberViewResponse
			memberView.UserID = member.UserID
			memberView.UserName = member.Name
			memberView.Avatar = member.Avatar
			memberView.Viewed = false
			for _, view := range *views {
				if view.ViewerID == member.UserID {
					memberView.Viewed = true
					memberView.ViewTime = view.Created
					if view.ViewerID == filter.UserID {
						(*list)[k].Viewed = true
					}
					break
				}
			}
			memberViews = append(memberViews, memberView)
		}
		(*list)[k].Views = memberViews
	}
	return list, err
}

func (s *projectService) GetProjectReportByID(reportID, userID, organizationID int64) (*ProjectReportResponse, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	memberQuery := member.NewMemberQuery(db)
	report, err := query.GetProjectReportByID(reportID, organizationID)
	if err != nil {
		msg := "报告不存在"
		return nil, errors.New(msg)
	}
	project, err := query.GetProjectByID(report.ProjectID, organizationID)
	if err != nil {
		msg := "项目不存在"
		return nil, errors.New(msg)
	}
	members, err := memberQuery.GetMembersByProjectID(project.ID)
	if err != nil {
		msg := "获取成员失败" + err.Error()
		return nil, errors.New(msg)
	}
	memberValid := false
	for _, member := range *members {
		if member.UserID == userID {
			memberValid = true
			break
		}
	}
	if organizationID == 0 {
		memberValid = true
	}
	if !memberValid {
		msg := "你不是此项目的成员"
		return nil, errors.New(msg)
	}
	links, err := query.GetProjectReportLinks(reportID)
	if err != nil {
		msg := "获取报告链接失败"
		return nil, errors.New(msg)
	}
	report.Links = *links
	return report, err
}

func (s *projectService) DeleteProjectReport(reportID, userID int64, userName string, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	// memberRepo := member.NewMemberRepository(tx)
	report, err := repo.GetProjectReportByID(reportID, organizationID)
	if err != nil {
		msg := "报告不存在"
		return errors.New(msg)
	}
	if report.UserID != userID {
		msg := "只能删除自己的报告"
		return errors.New(msg)
	}
	// project, err := repo.GetProjectByID(report.ProjectID, organizationID)
	// if err != nil {
	// 	msg := "项目不存在"
	// 	return errors.New(msg)
	// }
	// members, err := memberRepo.GetMembersByProjectID(project.ID)
	// if err != nil {
	// 	msg := "获取成员失败" + err.Error()
	// 	return errors.New(msg)
	// }
	// memberValid := false
	// for _, member := range *members {
	// 	if member.UserID == userID {
	// 		memberValid = true
	// 		break
	// 	}
	// }
	// if !memberValid {
	// 	msg := "你不是此项目的成员"
	// 	return errors.New(msg)
	// }
	err = repo.DeleteProjectReport(reportID, userName)
	if err != nil {
		msg := "删除报告失败" + err.Error()
		return errors.New(msg)
	}
	err = repo.DeleteProjectReportLinks(reportID, userName)
	if err != nil {
		msg := "删除报告链接失败" + err.Error()
		return errors.New(msg)
	}
	err = repo.DeleteProjectReportViews(reportID, userName)
	if err != nil {
		msg := "删除报告阅读记录失败" + err.Error()
		return errors.New(msg)
	}
	tx.Commit()
	return err
}

func (s *projectService) UpdateProjectReport(reportID int64, info ProjectReportNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	// memberRepo := member.NewMemberRepository(tx)
	oldReport, err := repo.GetProjectReportByID(reportID, info.OrganizationID)
	if err != nil {
		msg := "报告不存在"
		return errors.New(msg)
	}
	if oldReport.UserID != info.UserID {
		msg := "只能更新自己创建的报告"
		return errors.New(msg)
	}
	// _, err = repo.GetProjectByID(oldReport.ProjectID, info.OrganizationID)
	// if err != nil {
	// 	msg := "项目不存在"
	// 	return errors.New(msg)
	// }
	// members, err := memberRepo.GetMembersByProjectID(oldReport.ProjectID)
	// if err != nil {
	// 	msg := "获取项目成员失败"
	// 	return errors.New(msg)
	// }
	// memberValid := false
	// for _, member := range *members {
	// 	if member.UserID == info.UserID {
	// 		memberValid = true
	// 		break
	// 	}
	// }
	// if !memberValid {
	// 	msg := "你不是此项目的成员"
	// 	return errors.New(msg)
	// }
	var newReport ProjectReport
	newReport.Content = info.Content
	newReport.ReportDate = info.ReportDate
	newReport.Name = info.Name
	newReport.Content = info.Content
	newReport.Status = 1
	newReport.Updated = time.Now()
	newReport.UpdatedBy = info.User
	err = repo.UpdateProjectReport(reportID, newReport)
	if err != nil {
		msg := "更新报告失败"
		return errors.New(msg)
	}
	err = repo.DeleteProjectReportLinks(reportID, info.User)
	if err != nil {
		msg := "更新报告失败"
		return errors.New(msg)
	}
	err = repo.DeleteProjectReportViews(reportID, info.User)
	if err != nil {
		msg := "更新报告失败"
		return errors.New(msg)
	}
	for _, link := range info.Links {
		var reportLink ProjectReportLink
		reportLink.OrganizationID = info.OrganizationID
		reportLink.ProjectID = oldReport.ProjectID
		reportLink.ProjectReportID = reportID
		reportLink.Link = link
		reportLink.Status = 1
		reportLink.Created = time.Now()
		reportLink.CreatedBy = info.User
		reportLink.Updated = time.Now()
		reportLink.UpdatedBy = info.User
		err = repo.CreateProjectReportLink(reportLink)
		if err != nil {
			msg := "创建链接失败"
			return errors.New(msg)
		}
	}
	tx.Commit()

	type NewProjectReportCreated struct {
		ProjectReportID int64 `json:"project_report_id"`
	}
	var newEvent NewProjectReportCreated
	newEvent.ProjectReportID = oldReport.ID
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewProjectReportCreated", msg)
	if err != nil {
		msg := "发布消息（NewProjectReportCreated）失败"
		return errors.New(msg)
	}
	return err
}

func (s *projectService) NewProjectRecord(projectID int64, info ProjectRecordNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	memberRepo := member.NewMemberRepository(tx)
	project, err := repo.GetProjectByID(projectID, info.OrganizationID)
	if err != nil {
		msg := "项目不存在"
		return errors.New(msg)
	}
	members, err := memberRepo.GetMembersByProjectID(projectID)
	if err != nil {
		msg := "获取项目成员失败"
		return errors.New(msg)
	}
	memberValid := false
	for _, member := range *members {
		if member.UserID == info.UserID {
			memberValid = true
			break
		}
	}
	if !memberValid {
		msg := "你不是此项目的成员"
		return errors.New(msg)
	}
	var newRecord ProjectRecord
	newRecord.OrganizationID = info.OrganizationID
	newRecord.ProjectID = projectID
	newRecord.ClientID = project.ClientID
	newRecord.UserID = info.UserID
	newRecord.Content = info.Content
	newRecord.Plan = info.Plan
	newRecord.RecordDate = info.RecordDate
	newRecord.Name = info.Name
	newRecord.Content = info.Content
	newRecord.Status = 1
	newRecord.Created = time.Now()
	newRecord.CreatedBy = info.User
	newRecord.Updated = time.Now()
	newRecord.UpdatedBy = info.User
	recordID, err := repo.CreateProjectRecord(newRecord)
	if err != nil {
		msg := "创建报告失败"
		return errors.New(msg)
	}
	for _, link := range info.Photos {
		var recordPhoto ProjectRecordPhoto
		recordPhoto.OrganizationID = info.OrganizationID
		recordPhoto.ProjectID = projectID
		recordPhoto.ProjectRecordID = recordID
		recordPhoto.Link = link
		recordPhoto.Status = 1
		recordPhoto.Created = time.Now()
		recordPhoto.CreatedBy = info.User
		recordPhoto.Updated = time.Now()
		recordPhoto.UpdatedBy = info.User
		err = repo.CreateProjectRecordPhoto(recordPhoto)
		if err != nil {
			fmt.Println(err.Error())
			msg := "创建图片失败"
			return errors.New(msg)
		}
	}
	err = repo.UpdateProjectRecordDate(projectID)
	if err != nil {
		msg := "更新最后报告日期失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *projectService) GetProjectRecordList(projectID int64, filter ProjectRecordFilter, userType int) (int, *[]ProjectRecordResponse, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	memberQuery := member.NewMemberQuery(db)
	_, err := query.GetProjectByID(projectID, filter.OrganizationID)
	if err != nil {
		msg := "项目不存在"
		return 0, nil, errors.New(msg)
	}
	projectClientUserID, err := query.GetProjectClientUserID(projectID)
	if err != nil {
		msg := "获取客户失败"
		return 0, nil, errors.New(msg)
	}
	if filter.UserID != projectClientUserID && filter.OrganizationID != 0 && userType != 1 {
		members, err := memberQuery.GetMembersByProjectID(projectID)
		if err != nil {
			msg := "获取成员失败" + err.Error()
			return 0, nil, errors.New(msg)
		}
		memberValid := false
		for _, member := range *members {
			if member.UserID == filter.UserID {
				memberValid = true
				break
			}
		}
		if !memberValid {
			msg := "你不是此项目的成员"
			return 0, nil, errors.New(msg)
		}
	}
	count, err := query.GetProjectRecordCount(projectID)
	if err != nil {
		msg := "获取记录数量失败" + err.Error()
		return 0, nil, errors.New(msg)
	}
	list, err := query.GetProjectRecordList(projectID, filter)
	if err != nil {
		msg := "获取记录失败" + err.Error()
		return 0, nil, errors.New(msg)
	}
	for k, v := range *list {
		photos, err := query.GetProjectRecordPhotos(v.ID)
		if err != nil {
			msg := "获取图片失败" + err.Error()
			return 0, nil, errors.New(msg)
		}
		(*list)[k].Photos = *photos
	}
	return count, list, err
}

func (s *projectService) GetProjectRecordByID(recordID, userID, organizationID int64, userType int) (*ProjectRecordResponse, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	memberQuery := member.NewMemberQuery(db)
	record, err := query.GetProjectRecordByID(recordID, organizationID)
	if err != nil {
		msg := "报告不存在"
		return nil, errors.New(msg)
	}
	project, err := query.GetProjectByID(record.ProjectID, organizationID)
	if err != nil {
		msg := "项目不存在"
		return nil, errors.New(msg)
	}
	projectClientUserID, err := query.GetProjectClientUserID(project.ID)
	if err != nil {
		msg := "获取客户失败"
		return nil, errors.New(msg)
	}
	if userID != projectClientUserID && organizationID != 0 && userType != 1 {
		members, err := memberQuery.GetMembersByProjectID(project.ID)
		if err != nil {
			msg := "获取成员失败" + err.Error()
			return nil, errors.New(msg)
		}
		memberValid := false
		for _, member := range *members {
			if member.UserID == userID {
				memberValid = true
				break
			}
		}
		if !memberValid {
			msg := "你不是此项目的成员"
			return nil, errors.New(msg)
		}
	}
	photos, err := query.GetProjectRecordPhotos(recordID)
	if err != nil {
		msg := "获取记录图片失败"
		return nil, errors.New(msg)
	}
	record.Photos = *photos
	return record, err
}

func (s *projectService) DeleteProjectRecord(recordID, userID int64, userName string, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	// memberRepo := member.NewMemberRepository(tx)
	record, err := repo.GetProjectRecordByID(recordID, organizationID)
	if err != nil {
		msg := "记录不存在"
		return errors.New(msg)
	}
	if record.UserID != userID {
		msg := "只能删除自己的记录"
		return errors.New(msg)
	}
	// project, err := repo.GetProjectByID(record.ProjectID, organizationID)
	// if err != nil {
	// 	msg := "项目不存在"
	// 	return errors.New(msg)
	// }
	// members, err := memberRepo.GetMembersByProjectID(project.ID)
	// if err != nil {
	// 	msg := "获取成员失败" + err.Error()
	// 	return errors.New(msg)
	// }
	// memberValid := false
	// for _, member := range *members {
	// 	if member.UserID == userID {
	// 		memberValid = true
	// 		break
	// 	}
	// }
	// if !memberValid {
	// 	msg := "你不是此项目的成员"
	// 	return errors.New(msg)
	// }
	err = repo.DeleteProjectRecord(recordID, userName)
	if err != nil {
		msg := "删除报告失败" + err.Error()
		return errors.New(msg)
	}
	err = repo.DeleteProjectRecordPhotos(recordID, userName)
	if err != nil {
		msg := "删除报告图片失败" + err.Error()
		return errors.New(msg)
	}
	err = repo.UpdateProjectRecordDate(record.ProjectID)
	if err != nil {
		msg := "更新最后报告日期失败"
		return errors.New(msg)
	}
	tx.Commit()
	return err
}

func (s *projectService) UpdateProjectRecord(recordID int64, info ProjectRecordNew, userType int) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	// memberRepo := member.NewMemberRepository(tx)
	oldRecord, err := repo.GetProjectRecordByID(recordID, info.OrganizationID)
	if err != nil {
		msg := "报告不存在"
		return errors.New(msg)
	}
	if userType != 1 && oldRecord.UserID != info.UserID {
		msg := "只能更新自己创建的报告"
		return errors.New(msg)
	}
	var newRecord ProjectRecord
	newRecord.Content = info.Content
	newRecord.Plan = info.Plan
	newRecord.RecordDate = info.RecordDate
	newRecord.Name = info.Name
	newRecord.Content = info.Content
	newRecord.Status = 1
	newRecord.Updated = time.Now()
	newRecord.UpdatedBy = info.User
	err = repo.UpdateProjectRecord(recordID, newRecord)
	if err != nil {
		msg := "更新报告失败"
		return errors.New(msg)
	}
	err = repo.DeleteProjectRecordPhotos(recordID, info.User)
	if err != nil {
		msg := "更新报告失败"
		return errors.New(msg)
	}
	for _, photo := range info.Photos {
		var recordPhoto ProjectRecordPhoto
		recordPhoto.OrganizationID = info.OrganizationID
		recordPhoto.ProjectID = oldRecord.ProjectID
		recordPhoto.ProjectRecordID = recordID
		recordPhoto.Link = photo
		recordPhoto.Status = 1
		recordPhoto.Created = time.Now()
		recordPhoto.CreatedBy = info.User
		recordPhoto.Updated = time.Now()
		recordPhoto.UpdatedBy = info.User
		err = repo.CreateProjectRecordPhoto(recordPhoto)
		if err != nil {
			msg := "创建图片失败"
			return errors.New(msg)
		}
	}
	err = repo.UpdateProjectRecordDate(oldRecord.ProjectID)
	if err != nil {
		msg := "更新最后报告日期失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *projectService) PortalGetProjectRecordList(projectID int64, filter ProjectRecordFilter) (int, *[]ProjectRecordResponse, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	_, err := query.GetProjectByID(projectID, filter.OrganizationID)
	if err != nil {
		msg := "项目不存在"
		return 0, nil, errors.New(msg)
	}
	count, err := query.GetProjectRecordCount(projectID)
	if err != nil {
		msg := "获取记录数量失败" + err.Error()
		return 0, nil, errors.New(msg)
	}
	list, err := query.GetProjectRecordList(projectID, filter)
	if err != nil {
		msg := "获取记录失败" + err.Error()
		return 0, nil, errors.New(msg)
	}
	for k, v := range *list {
		photos, err := query.GetProjectRecordPhotos(v.ID)
		if err != nil {
			msg := "获取图片失败" + err.Error()
			return 0, nil, errors.New(msg)
		}
		(*list)[k].Photos = *photos
	}
	return count, list, err
}

func (s *projectService) ViewProjectReport(reportID, organizationID, userID int64, userName string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewProjectRepository(tx)
	memberRepo := member.NewMemberRepository(tx)
	oldReport, err := repo.GetProjectReportByID(reportID, organizationID)
	if err != nil {
		msg := "报告不存在"
		return errors.New(msg)
	}
	_, err = repo.GetProjectByID(oldReport.ProjectID, organizationID)
	if err != nil {
		msg := "项目不存在"
		return errors.New(msg)
	}
	members, err := memberRepo.GetMembersByProjectID(oldReport.ProjectID)
	if err != nil {
		msg := "获取项目成员失败"
		return errors.New(msg)
	}
	memberValid := false
	for _, member := range *members {
		if member.UserID == userID {
			memberValid = true
			break
		}
	}
	if !memberValid {
		msg := "你不是此项目的成员"
		return errors.New(msg)
	}
	count, err := repo.CheckViewExist(reportID, userID)
	if err != nil {
		msg := "获取已阅记录失败"
		return errors.New(msg)
	}
	if count != 0 {
		msg := "重复确认"
		return errors.New(msg)
	}
	var newReportView ProjectReportView
	newReportView.OrganizationID = organizationID
	newReportView.ProjectID = oldReport.ProjectID
	newReportView.ProjectReportID = reportID
	newReportView.ViewerID = userID
	newReportView.ViewerName = userName
	newReportView.Status = 1
	newReportView.Created = time.Now()
	newReportView.CreatedBy = userName
	newReportView.Updated = time.Now()
	newReportView.UpdatedBy = userName
	err = repo.CreateProjectReportView(newReportView)
	if err != nil {
		msg := "创建已阅失败"
		return errors.New(msg)
	}
	views, err := repo.GetProjectReportView(reportID)
	if err != nil {
		msg := "获取阅读记录失败"
		return errors.New(msg)
	}
	reportStatus := 1
	if len(*views) == len(*members) {
		reportStatus = 3
	} else {
		reportStatus = 2
	}
	err = repo.UpdateProjectReportStatus(reportID, reportStatus, userName)
	if err != nil {
		msg := "更新状态失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *projectService) GetProjectReportUnreadList(userID int64) (*[]ProjectReportResponse, error) {
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	list, err := query.GetProjectReportUnreadList(userID)
	if err != nil {
		fmt.Println(err.Error())
		msg := "获取未读报告失败" // + err.Error()
		return nil, errors.New(msg)
	}
	return list, err
}

func (s *projectService) GetProjectRecordStatus(projectID, organizationID int64) (*ProjectRecordStatusResponse, error) {
	var res ProjectRecordStatusResponse
	var lastRecordDate string
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	project, err := query.GetProjectByID(projectID, organizationID)
	if err != nil {
		msg := "项目不存在"
		return nil, errors.New(msg)
	}
	count, err := query.GetProjectRecordCount(projectID)
	if err != nil {
		msg := "获取项目记录数量失败"
		return nil, errors.New(msg)
	}
	res.StartDate = project.Created.Format("2006-01-02")
	res.RecordCount = count
	res.LastRecordDate = project.LastRecordDate
	if project.LastRecordDate == "" {
		lastRecordDate = res.StartDate
	} else {
		lastRecordDate = project.LastRecordDate
	}
	recordDateTime, err := time.Parse("2006-01-02", lastRecordDate)
	if err != nil {
		return nil, err
	}
	today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	interval := today.Sub(recordDateTime)
	res.NoRecordDay = int(interval.Hours() / 24)
	return &res, nil

}

func (s *projectService) GetProjectSumByStatus(filter ProjectSumFilter, organizationID int64) (*[]ProjectSumByStatus, error) {
	if organizationID == 0 && filter.OrganizationID == 0 {
		msg := "组织ID不能为空"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	res, err := query.GetProjectSumByStatus(filter)
	return res, err
}

func (s *projectService) GetProjectSumByTeam(filter ProjectSumFilter, organizationID int64) (*[]ProjectSumByTeam, error) {
	if organizationID == 0 && filter.OrganizationID == 0 {
		msg := "组织ID不能为空"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	res, err := query.GetProjectSumByTeam(filter)
	return res, err
}

func (s *projectService) GetProjectSumByUser(filter ProjectSumFilter, organizationID int64) (*[]ProjectSumByUser, error) {
	if organizationID == 0 && filter.OrganizationID == 0 {
		msg := "组织ID不能为空"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	res, err := query.GetProjectSumByUser(filter)
	return res, err
}

func (s *projectService) GetProjectSumByArea(filter ProjectSumFilter, organizationID int64) (*[]ProjectSumByArea, error) {
	if organizationID == 0 && filter.OrganizationID == 0 {
		msg := "组织ID不能为空"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewProjectQuery(db)
	res, err := query.GetProjectSumByArea(filter)
	return res, err
}
