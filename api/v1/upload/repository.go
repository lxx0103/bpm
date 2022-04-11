package upload

import (
	"database/sql"
	"time"
)

type uploadRepository struct {
	tx *sql.Tx
}

func NewUploadRepository(transaction *sql.Tx) UploadRepository {
	return &uploadRepository{
		tx: transaction,
	}
}

type UploadRepository interface {
	//Upload Management
	CreateUpload(string, int64, string) error
}

func (r *uploadRepository) CreateUpload(path string, organizationID int64, user string) error {
	_, err := r.tx.Exec(`
			INSERT INTO file_uploads
			(
				path,
				organization_id,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, path, organizationID, 1, time.Now(), user, time.Now(), user)
	return err
}
