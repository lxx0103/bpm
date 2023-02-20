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
	tx.Commit()

	type NewProjectReportCreated struct {
		ProjectID int64 `json:"project_id"`
	}
	var newEvent NewProjectReportCreated
	newEvent.ProjectID = projectID
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
		ProjectID int64 `json:"project_id"`
	}
	var newEvent NewProjectReportCreated
	newEvent.ProjectID = oldReport.ProjectID
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
			msg := "创建图片失败"
			return errors.New(msg)
		}
	}
	tx.Commit()
	return nil
}

func (s *projectService) GetProjectRecordList(projectID int64, filter ProjectRecordFilter) (int, *[]ProjectRecordResponse, error) {
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
	if filter.UserID != projectClientUserID && filter.OrganizationID != 0 {
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

func (s *projectService) GetProjectRecordByID(recordID, userID, organizationID int64) (*ProjectRecordResponse, error) {
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
	if userID != projectClientUserID && organizationID != 0 {
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
	tx.Commit()
	return err
}

func (s *projectService) UpdateProjectRecord(recordID int64, info ProjectRecordNew) error {
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
	if oldRecord.UserID != info.UserID {
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
