package http

import (
	"log"
	"net/http"
	"quickshare/core/model"
	"quickshare/core/service"
	web "quickshare/pkg"
	"time"

	"github.com/gorilla/mux"
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

func (h *UploadObjectHandler) ConfirmUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing upload id"})
		return
	}

	log.Println("confirming upload for id", id)

	if err := h.uploadObjectService.ConfirmUpload(id); err != nil {
		log.Println("error confirming upload", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	web.WriteJSON(w, http.StatusOK, map[string]string{
		"id":     id,
		"status": "completed",
	})
}

// Download retorna uma URL de download para um upload já concluído
func (h *UploadObjectHandler) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing upload id"})
		return
	}

	log.Println("generating download URL for id", id)

	downloadURL, err := h.uploadObjectService.GetDownloadURL(id)
	if err != nil {
		log.Println("error generating download URL", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	web.WriteJSON(w, http.StatusOK, map[string]string{
		"id":           id,
		"download_url": downloadURL,
	})
}