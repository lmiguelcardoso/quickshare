package repository

import (
	"database/sql"
	"quickshare/core/model"
)

type PostgreSQLRepository struct {
	db *sql.DB
}
func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) CreateUploadObject(uploadObject *model.UploadObject) (*model.UploadObject, error) {
	query := `INSERT INTO upload_objects (id, file_name, file_size, mime_type, object_key, status, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(query, uploadObject.ID, uploadObject.FileName, uploadObject.FileSize, uploadObject.MimeType, uploadObject.ObjectKey, uploadObject.Status, uploadObject.ExpiresAt).Scan(&uploadObject.ID)
	if err != nil {
		return nil, err
	}
	return uploadObject, nil
}

func (r *PostgreSQLRepository) GetUploadObject(id string) (*model.UploadObject, error) {
	query := `SELECT id, file_name, file_size, mime_type, object_key, status, expires_at FROM upload_objects WHERE id = $1`
	var uploadObject model.UploadObject
	
	row := r.db.QueryRow(query, id)
	err := row.Scan(&uploadObject.ID, &uploadObject.FileName, &uploadObject.FileSize, &uploadObject.MimeType, &uploadObject.ObjectKey, &uploadObject.Status, &uploadObject.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &uploadObject, nil
}

func (r *PostgreSQLRepository) UpdateUploadObject(id string, uploadObject *model.UploadObject) (*model.UploadObject, error) {
	query := `UPDATE upload_objects SET file_name = $1, file_size = $2, mime_type = $3, object_key = $4, status = $5, expires_at = $6 WHERE id = $7 RETURNING id`
	err := r.db.QueryRow(query, uploadObject.FileName, uploadObject.FileSize, uploadObject.MimeType, uploadObject.ObjectKey, uploadObject.Status, uploadObject.ExpiresAt, id).Scan(&uploadObject.ID)
	if err != nil {
		return nil, err
	}
	return uploadObject, nil
}

func (r *PostgreSQLRepository) DeleteUploadObject(id string) error {
	query := `DELETE FROM upload_objects WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}