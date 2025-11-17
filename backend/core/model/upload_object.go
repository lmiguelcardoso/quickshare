package model

import "time"

type UploadObject struct {
	ID string `json:"id"`
	FileName string `json:"file_name"`
	FileSize int64 `json:"file_size"`
	MimeType string `json:"mime_type"`
	ObjectKey string `json:"object_key"`
	Status string `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}