package service

import (
	"errors"
	"fmt"
	"quickshare/core/model"
	"quickshare/core/repository"
	"time"
)
type UploadResponse struct {
	ID        string    `json:"id"`
	UploadURL string    `json:"upload_url"`
	ObjectKey string    `json:"object_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UploadObjectService struct {
	repository  repository.UploadObjectRepository
	blobStorage repository.BlobStorageRepository
}

func NewUploadObjectService(repo repository.UploadObjectRepository, blobStorage repository.BlobStorageRepository) *UploadObjectService {
	return &UploadObjectService{
		repository:  repo,
		blobStorage: blobStorage,
	}
}

func (s *UploadObjectService) generateID(uploadObject *model.UploadObject) string {
	fileName := uploadObject.FileName[:2]
	fileSize := fmt.Sprintf("%d", uploadObject.FileSize)[:2]
	now := fmt.Sprintf("%d", time.Now().Unix())[:2]

	return fmt.Sprintf("%s%s%s", fileName, fileSize, now)
}

func (s *UploadObjectService) DeleteUploadObject(id string) error {
	uploadObject, err := s.repository.GetUploadObject(id)
	if err != nil {
		return err
	}

	if err := s.blobStorage.Delete(uploadObject.ID); err != nil {
		return fmt.Errorf("failed to delete from storage: %w", err)
	}

	return s.repository.DeleteUploadObject(id)
}

func (s *UploadObjectService) InitiateUpload(uploadObject *model.UploadObject) (*UploadResponse, error) {
	// 1. create upload object
	uploadObject.ID = s.generateID(uploadObject)
	uploadObject.Status = "pending"
	uploadObject.ObjectKey = fmt.Sprintf("uploads/%s/%s", uploadObject.ID, uploadObject.FileName)

	if uploadObject.ExpiresAt.IsZero() {
		uploadObject.ExpiresAt = time.Now().Add(time.Hour * 24)
	}

	// 2. create presigned URL for upload (expires in 15 minutes)
	uploadURL, err := s.blobStorage.GeneratePresignedUploadURL(uploadObject.ObjectKey, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	// 3. save to database
	created, err := s.repository.CreateUploadObject(uploadObject)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload object: %w", err)
	}

	return &UploadResponse{
		ID:         created.ID,
		UploadURL:  uploadURL,
		ObjectKey:  created.ObjectKey,
		ExpiresAt:  created.ExpiresAt,
	}, nil
}

// ConfirmUpload check if file was uploaded and update status
func (s *UploadObjectService) ConfirmUpload(id string) error {
	// 1. get upload object from database
	uploadObject, err := s.repository.GetUploadObject(id)
	if err != nil {
		return err
	}

	// 2. check if object exists in storage
	exists, err := s.blobStorage.ObjectExists(uploadObject.ObjectKey)
	if err != nil {
		return fmt.Errorf("failed to verify object: %w", err)
	}

	if !exists {
		return errors.New("file not found in storage")
	}

	// 3. get object metadata
	metadata, err := s.blobStorage.GetObjectMetadata(uploadObject.ObjectKey)
	if err == nil {
		// update file size if available
		if size, ok := metadata["size"]; ok {
			fmt.Sscanf(size, "%d", &uploadObject.FileSize)
		}
	}

		// 4. Atualiza status para "completed"
	uploadObject.Status = "completed"
	_, err = s.repository.UpdateUploadObject(id, uploadObject)
	return err
}

// GetDownloadURL create presigned URL for download
func (s *UploadObjectService) GetDownloadURL(id string) (string, error) {
	uploadObject, err := s.repository.GetUploadObject(id)
	if err != nil {
		return "", err
	}

	if uploadObject.Status != "completed" {
		return "", errors.New("upload not completed yet")
	}

	if time.Now().After(uploadObject.ExpiresAt) {
		return "", errors.New("upload has expired")
	}

    return s.blobStorage.GetPublicURL(uploadObject.ObjectKey), nil
}
