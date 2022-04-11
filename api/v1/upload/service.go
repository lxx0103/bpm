package upload

import (
	"bpm/core/database"
)

type uploadService struct {
}

func NewUploadService() UploadService {
	return &uploadService{}
}

// UploadService represents a service for managing uploads.
type UploadService interface {
	//Upload Management
	NewUpload(string, int64, string) error
	GetUploadList(UploadFilter, int64) (int, *[]Upload, error)
}

func (s *uploadService) NewUpload(path string, organizationID int64, userName string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewUploadRepository(tx)
	err = repo.CreateUpload(path, organizationID, userName)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (s *uploadService) GetUploadList(filter UploadFilter, organizationID int64) (int, *[]Upload, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewUploadQuery(db)
	count, err := query.GetUploadCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetUploadList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, nil
}
