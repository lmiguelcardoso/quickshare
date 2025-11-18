package repository

import (
	"quickshare/core/model"
)

type UploadObjectRepository interface {
	CreateUploadObject(uploadObject *model.UploadObject) (*model.UploadObject, error)
	GetUploadObject(id string) (*model.UploadObject, error)
	UpdateUploadObjectStatus(id string, status string) error
	DeleteUploadObject(id string) error
}
