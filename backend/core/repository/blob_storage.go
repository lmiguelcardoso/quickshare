package repository

import "time"

type BlobStorageRepository interface {
	GeneratePresignedUploadURL(objectKey string, expiresIn time.Duration) (string, error)
	GeneratePresignedDownloadURL(objectKey string, expiresIn time.Duration) (string, error)
	ObjectExists(objectKey string) (bool, error)
	GetObjectMetadata(objectKey string) (map[string]string, error)
}