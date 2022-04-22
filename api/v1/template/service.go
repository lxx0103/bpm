package template

import (
	"bpm/core/database"
	"errors"
)

type templateService struct {
}

func NewTemplateService() TemplateService {
	return &templateService{}
}

// TemplateService represents a service for managing templates.
type TemplateService interface {
	//Template Management
	GetTemplateByID(int64, int64) (*Template, error)
	NewTemplate(TemplateNew, int64) (*Template, error)
	GetTemplateList(TemplateFilter, int64) (int, *[]Template, error)
	UpdateTemplate(int64, TemplateUpdate, int64) (*Template, error)
	DeleteTemplate(int64, int64, string) error
}

func (s *templateService) GetTemplateByID(id int64, organizationID int64) (*Template, error) {
	db := database.InitMySQL()
	query := NewTemplateQuery(db)
	template, err := query.GetTemplateByID(id, organizationID)
	return template, err
}

func (s *templateService) NewTemplate(info TemplateNew, organizationID int64) (*Template, error) {
	if organizationID != 0 && organizationID != info.OrganizationID {
		msg := "无权新建模板"
		return nil, errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewTemplateRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, 0)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "模板名称重复"
		return nil, errors.New(msg)
	}
	templateID, err := repo.CreateTemplate(info)
	if err != nil {
		return nil, err
	}
	template, err := repo.GetTemplateByID(templateID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return template, err
}

func (s *templateService) GetTemplateList(filter TemplateFilter, organizationID int64) (int, *[]Template, error) {
	db := database.InitMySQL()
	query := NewTemplateQuery(db)
	count, err := query.GetTemplateCount(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetTemplateList(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *templateService) UpdateTemplate(templateID int64, info TemplateUpdate, organizationID int64) (*Template, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewTemplateRepository(tx)
	oldTemplate, err := repo.GetTemplateByID(templateID)
	if err != nil {
		return nil, err
	}
	if organizationID != 0 && organizationID != oldTemplate.OrganizationID {
		msg := "你无权修改此模板"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, organizationID, templateID)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "模板名称重复"
		return nil, errors.New(msg)
	}
	if info.Name != "" {
		oldTemplate.Name = info.Name
	}
	if info.Type != 0 {
		oldTemplate.Type = info.Type
	}
	if info.Status != 0 {
		oldTemplate.Status = info.Status
	}
	if info.EventJson != "" {
		oldTemplate.EventJson = info.EventJson
	}
	err = repo.UpdateTemplate(templateID, *oldTemplate, info.User)
	if err != nil {
		return nil, err
	}
	template, err := repo.GetTemplateByID(templateID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return template, err
}
func (s *templateService) DeleteTemplate(templateID int64, organizationID int64, user string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewTemplateRepository(tx)
	oldTemplate, err := repo.GetTemplateByID(templateID)
	if err != nil {
		return err
	}
	if organizationID != 0 && organizationID != oldTemplate.OrganizationID {
		msg := "你无权删除此模板"
		return errors.New(msg)
	}
	err = repo.DeleteTemplate(templateID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
