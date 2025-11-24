package http

import (
	"log"
	"net/http"
	"quickshare/core/model"
	"quickshare/core/service"
	web "quickshare/pkg"
	"time"
)

type uploadObjectRequest struct {
	FileName string `json:"file_name"`
	FileSize int64 `json:"file_size"`
	MimeType string `json:"mime_type"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UploadObjectHandler struct {
	uploadObjectService *service.UploadObjectService
}

func NewUploadObjectHandler(uploadObjectService *service.UploadObjectService) *UploadObjectHandler {
	return &UploadObjectHandler{uploadObjectService: uploadObjectService}
}

func (h *UploadObjectHandler) UploadObject(w http.ResponseWriter, r *http.Request) {
	var uploadObject model.UploadObject
	
	log.Println("starting upload object", uploadObject.FileName)

	if err := web.ReadJSON(r, &uploadObject); err != nil {
		log.Println("error reading request body", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	uploadResponse, err := h.uploadObjectService.InitiateUpload(&uploadObject)
	if err != nil {
		log.Println("error initiating upload", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to initiate upload"})
		return
	}

	web.WriteJSON(w, http.StatusOK, uploadResponse)
}