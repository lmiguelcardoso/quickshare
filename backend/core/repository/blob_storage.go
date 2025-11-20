package repository

import "time"

type BlobStorageRepository interface {
	GeneratePresignedUploadURL(objectKey string, expiresIn time.Duration) (string, error)
	GetPublicURL(objectKey string) string
	ObjectExists(objectKey string) (bool, error)
	GetObjectMetadata(objectKey string) (map[string]string, error)
	Delete(objectKey string) error
}