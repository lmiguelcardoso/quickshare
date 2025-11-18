package repository

import (
	"quickshare/core/model"
)

type UploadObjectRepository interface {
	CreateUploadObject(uploadObject *model.UploadObject) (*model.UploadObject, error)
	GetUploadObject(id string) (*model.UploadObject, error)
	UpdateUploadObject(id string, uploadObject *model.UploadObject) (*model.UploadObject, error)
	DeleteUploadObject(id string) error
}
