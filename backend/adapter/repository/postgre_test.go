package repository

import (
	"database/sql"
	"errors"
	"quickshare/core/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgreSQLRepository_CreateUploadObject(t *testing.T) {
	fixedTime := time.Date(2024, 11, 18, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		input       *model.UploadObject
		mockSetup   func(mock sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - create upload object",
			input: &model.UploadObject{
				ID:        "test-id-123",
				FileName:  "document.pdf",
				FileSize:  1024,
				MimeType:  "application/pdf",
				ObjectKey: "uploads/document.pdf",
				Status:    "pending",
				ExpiresAt: fixedTime,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("test-id-123")
				mock.ExpectQuery(`INSERT INTO upload_objects`).
					WithArgs("test-id-123", "document.pdf", int64(1024), "application/pdf", "uploads/document.pdf", "pending", fixedTime).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error - duplicate id",
			input: &model.UploadObject{
				ID:        "duplicate-id",
				FileName:  "test.txt",
				FileSize:  512,
				MimeType:  "text/plain",
				ObjectKey: "uploads/test.txt",
				Status:    "pending",
				ExpiresAt: fixedTime,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO upload_objects`).
					WithArgs("duplicate-id", "test.txt", int64(512), "text/plain", "uploads/test.txt", "pending", fixedTime).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			wantErr:     true,
			errContains: "duplicate key",
		},
		{
			name: "error - database connection failure",
			input: &model.UploadObject{
				ID:        "test-id-456",
				FileName:  "image.jpg",
				FileSize:  2048,
				MimeType:  "image/jpeg",
				ObjectKey: "uploads/image.jpg",
				Status:    "pending",
				ExpiresAt: fixedTime,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO upload_objects`).
					WithArgs("test-id-456", "image.jpg", int64(2048), "image/jpeg", "uploads/image.jpg", "pending", fixedTime).
					WillReturnError(errors.New("connection refused"))
			},
			wantErr:     true,
			errContains: "connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewPostgreSQLRepository(db)
			result, err := repo.CreateUploadObject(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ID != tt.input.ID {
				t.Errorf("expected ID %q, got %q", tt.input.ID, result.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQLRepository_GetUploadObject(t *testing.T) {
	fixedTime := time.Date(2024, 11, 18, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		inputID     string
		mockSetup   func(mock sqlmock.Sqlmock)
		want        *model.UploadObject
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - get upload object",
			inputID: "test-id-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "file_name", "file_size", "mime_type", "object_key", "status", "expires_at"}).
					AddRow("test-id-123", "document.pdf", 1024, "application/pdf", "uploads/document.pdf", "active", fixedTime)
				mock.ExpectQuery(`SELECT id, file_name, file_size, mime_type, object_key, status, expires_at FROM upload_objects WHERE id`).
					WithArgs("test-id-123").
					WillReturnRows(rows)
			},
			want: &model.UploadObject{
				ID:        "test-id-123",
				FileName:  "document.pdf",
				FileSize:  1024,
				MimeType:  "application/pdf",
				ObjectKey: "uploads/document.pdf",
				Status:    "active",
				ExpiresAt: fixedTime,
			},
			wantErr: false,
		},
		{
			name:    "error - not found",
			inputID: "non-existent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, file_name, file_size, mime_type, object_key, status, expires_at FROM upload_objects WHERE id`).
					WithArgs("non-existent-id").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:     true,
			errContains: "no rows",
		},
		{
			name:    "error - database error",
			inputID: "error-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, file_name, file_size, mime_type, object_key, status, expires_at FROM upload_objects WHERE id`).
					WithArgs("error-id").
					WillReturnError(errors.New("database connection lost"))
			},
			wantErr:     true,
			errContains: "connection lost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewPostgreSQLRepository(db)
			result, err := repo.GetUploadObject(tt.inputID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ID != tt.want.ID {
				t.Errorf("expected ID %q, got %q", tt.want.ID, result.ID)
			}
			if result.FileName != tt.want.FileName {
				t.Errorf("expected FileName %q, got %q", tt.want.FileName, result.FileName)
			}
			if result.FileSize != tt.want.FileSize {
				t.Errorf("expected FileSize %d, got %d", tt.want.FileSize, result.FileSize)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQLRepository_UpdateUploadObject(t *testing.T) {
	fixedTime := time.Date(2024, 11, 18, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		inputID     string
		input       *model.UploadObject
		mockSetup   func(mock sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - update upload object",
			inputID: "test-id-123",
			input: &model.UploadObject{
				ID:        "test-id-123",
				FileName:  "updated-document.pdf",
				FileSize:  2048,
				MimeType:  "application/pdf",
				ObjectKey: "uploads/updated-document.pdf",
				Status:    "completed",
				ExpiresAt: fixedTime,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("test-id-123")
				mock.ExpectQuery(`UPDATE upload_objects SET`).
					WithArgs("updated-document.pdf", int64(2048), "application/pdf", "uploads/updated-document.pdf", "completed", fixedTime, "test-id-123").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:    "error - not found",
			inputID: "non-existent-id",
			input: &model.UploadObject{
				FileName:  "test.txt",
				FileSize:  512,
				MimeType:  "text/plain",
				ObjectKey: "uploads/test.txt",
				Status:    "pending",
				ExpiresAt: fixedTime,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE upload_objects SET`).
					WithArgs("test.txt", int64(512), "text/plain", "uploads/test.txt", "pending", fixedTime, "non-existent-id").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:     true,
			errContains: "no rows",
		},
		{
			name:    "error - database error",
			inputID: "error-id",
			input: &model.UploadObject{
				FileName:  "test.txt",
				FileSize:  512,
				MimeType:  "text/plain",
				ObjectKey: "uploads/test.txt",
				Status:    "pending",
				ExpiresAt: fixedTime,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE upload_objects SET`).
					WithArgs("test.txt", int64(512), "text/plain", "uploads/test.txt", "pending", fixedTime, "error-id").
					WillReturnError(errors.New("update failed"))
			},
			wantErr:     true,
			errContains: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewPostgreSQLRepository(db)
			result, err := repo.UpdateUploadObject(tt.inputID, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ID != tt.input.ID {
				t.Errorf("expected ID %q, got %q", tt.input.ID, result.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQLRepository_DeleteUploadObject(t *testing.T) {
	tests := []struct {
		name        string
		inputID     string
		mockSetup   func(mock sqlmock.Sqlmock)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - delete upload object",
			inputID: "test-id-123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM upload_objects WHERE id`).
					WithArgs("test-id-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:    "success - delete non-existent object (no error)",
			inputID: "non-existent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM upload_objects WHERE id`).
					WithArgs("non-existent-id").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
		{
			name:    "error - database error",
			inputID: "error-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM upload_objects WHERE id`).
					WithArgs("error-id").
					WillReturnError(errors.New("delete failed"))
			},
			wantErr:     true,
			errContains: "delete failed",
		},
		{
			name:    "error - connection timeout",
			inputID: "timeout-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM upload_objects WHERE id`).
					WithArgs("timeout-id").
					WillReturnError(errors.New("connection timeout"))
			},
			wantErr:     true,
			errContains: "timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewPostgreSQLRepository(db)
			err = repo.DeleteUploadObject(tt.inputID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}