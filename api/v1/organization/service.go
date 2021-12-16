package organization

import (
	"bpm/core/database"
)

type organizationService struct {
}

func NewOrganizationService() OrganizationService {
	return &organizationService{}
}

// OrganizationService represents a service for managing organizations.
type OrganizationService interface {
	//Organization Management
	GetOrganizationByID(int64) (*Organization, error)
	NewOrganization(OrganizationNew) (*Organization, error)
	GetOrganizationList(OrganizationFilter) (int, *[]Organization, error)
	UpdateOrganization(int64, OrganizationNew) (*Organization, error)
}

func (s *organizationService) GetOrganizationByID(id int64) (*Organization, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	organization, err := query.GetOrganizationByID(id)
	return organization, err
}

func (s *organizationService) NewOrganization(info OrganizationNew) (*Organization, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewOrganizationRepository(tx)
	organizationID, err := repo.CreateOrganization(info)
	if err != nil {
		return nil, err
	}
	organization, err := repo.GetOrganizationByID(organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return organization, err
}

func (s *organizationService) GetOrganizationList(filter OrganizationFilter) (int, *[]Organization, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	count, err := query.GetOrganizationCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetOrganizationList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *organizationService) UpdateOrganization(organizationID int64, info OrganizationNew) (*Organization, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewOrganizationRepository(tx)
	_, err = repo.UpdateOrganization(organizationID, info)
	if err != nil {
		return nil, err
	}
	organization, err := repo.GetOrganizationByID(organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return organization, err
}
