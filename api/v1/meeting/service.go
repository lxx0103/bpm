package meeting

import (
	"bpm/core/database"
	"errors"
)

type meetingService struct {
}

func NewMeetingService() *meetingService {
	return &meetingService{}
}

func (s *meetingService) GetMeetingByID(id int64, organizationID int64) (*MeetingResponse, error) {
	db := database.InitMySQL()
	query := NewMeetingQuery(db)
	meeting, err := query.GetMeetingByID(id, organizationID)
	return meeting, err
}

func (s *meetingService) NewMeeting(info MeetingNew, organizationID int64) error {
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
	repo := NewMeetingRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, 0)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "会议名称重复"
		return errors.New(msg)
	}
	err = repo.CreateMeeting(info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *meetingService) GetMeetingList(filter MeetingFilter, organizationID int64) (int, *[]MeetingResponse, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewMeetingQuery(db)
	count, err := query.GetMeetingCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetMeetingList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *meetingService) UpdateMeeting(meetingID int64, info MeetingNew, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewMeetingRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, meetingID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "会议记录名称重复"
		return errors.New(msg)
	}
	oldMeeting, err := repo.GetMeetingByID(meetingID)
	if err != nil {
		msg := "会议记录不存在"
		return errors.New(msg)
	}
	if organizationID != oldMeeting.OrganizationID && organizationID != 0 {
		msg := "会议记录不存在"
		return errors.New(msg)
	}
	if oldMeeting.UserID != info.UserID {
		msg := "不是你创建的会议记录"
		return errors.New(msg)
	}
	err = repo.UpdateMeeting(meetingID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *meetingService) DeleteMeeting(meetingID, organizationID int64, byUser string, byUserID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewMeetingRepository(tx)
	oldMeeting, err := repo.GetMeetingByID(meetingID)
	if err != nil {
		msg := "会议记录不存在"
		return errors.New(msg)
	}
	if organizationID != oldMeeting.OrganizationID && organizationID != 0 {
		msg := "会议记录不存在"
		return errors.New(msg)
	}
	if oldMeeting.UserID != byUserID {
		msg := "不是你创建的会议记录"
		return errors.New(msg)
	}
	err = repo.DeleteMeeting(meetingID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
