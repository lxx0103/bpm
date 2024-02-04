package team

import (
	"bpm/core/database"
	"errors"
)

type teamService struct {
}

func NewTeamService() *teamService {
	return &teamService{}
}

func (s *teamService) GetTeamByID(id int64, organizationID int64) (*Team, error) {
	db := database.InitMySQL()
	query := NewTeamQuery(db)
	team, err := query.GetTeamByID(id, organizationID)
	return team, err
}

func (s *teamService) NewTeam(info TeamNew, organizationID int64) (*Team, error) {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewTeamRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, 0)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "班组名称重复"
		return nil, errors.New(msg)
	}
	teamID, err := repo.CreateTeam(info)
	if err != nil {
		return nil, err
	}
	team, err := repo.GetTeamByID(teamID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return team, err
}

func (s *teamService) GetTeamList(filter TeamFilter, organizationID int64) (int, *[]TeamResponse, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewTeamQuery(db)
	count, err := query.GetTeamCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetTeamList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *teamService) UpdateTeam(teamID int64, info TeamNew, organizationID int64) (*Team, error) {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewTeamRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, teamID)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "班组名称重复"
		return nil, errors.New(msg)
	}
	_, err = repo.UpdateTeam(teamID, info)
	if err != nil {
		return nil, err
	}
	team, err := repo.GetTeamByID(teamID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return team, err
}
