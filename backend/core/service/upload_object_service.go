package service

import "quickshare/core/model"

type UploadObjectService interface {
	CreateUploadObject(uploadObject *model.UploadObject) (*model.UploadObject, error)
	GetUploadObject(id string) (*model.UploadObject, error)
	UpdateUploadObject(id string, uploadObject *model.UploadObject) (*model.UploadObject, error)
	DeleteUploadObject(id string) error
}
