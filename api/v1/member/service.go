package member

import (
	"bpm/core/database"
	"bpm/core/queue"
	"encoding/json"
	"errors"
)

type memberService struct {
}

func NewMemberService() *memberService {
	return &memberService{}
}

func (s *memberService) NewMember(info MemberNew, organizationID int64) (*[]MemberResponse, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewMemberRepository(tx)
	projectExist, err := repo.CheckProjectExist(info.ProjectID, organizationID)
	if err != nil {
		return nil, err
	}
	if projectExist == 0 {
		msg := "项目不存在"
		return nil, errors.New(msg)
	}
	err = repo.DeleteProjectMember(info.ProjectID, info.User)
	if err != nil {
		return nil, err
	}
	err = repo.CreateProjectMember(info.ProjectID, info.UserID, organizationID, info.User)
	if err != nil {
		return nil, err
	}
	members, err := repo.GetMembersByProjectID(info.ProjectID)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	type NewProjectCreated struct {
		ProjectID int64 `json:"project_id"`
	}
	var newEvent NewProjectCreated
	newEvent.ProjectID = info.ProjectID
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewProjectMember", msg)
	if err != nil {
		msg := "create event NewProjectMember error"
		return nil, errors.New(msg)
	}
	return members, err
}

func (s *memberService) GetMemberList(projectID int64, organizationID int64) (*[]MemberResponse, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewMemberRepository(tx)
	projectExist, err := repo.CheckProjectExist(projectID, organizationID)
	if err != nil {
		return nil, err
	}
	if projectExist == 0 {
		msg := "项目不存在"
		return nil, errors.New(msg)
	}
	members, err := repo.GetMembersByProjectID(projectID)
	if err != nil {
		return nil, err
	}
	return members, err
}
