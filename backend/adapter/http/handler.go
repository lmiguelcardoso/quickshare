package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
	uploadObjectHandler *UploadObjectHandler
}

func NewHandler(uploadObjectHandler *UploadObjectHandler) *handler {
	return &handler{
		uploadObjectHandler: uploadObjectHandler,
	}
}

func (h *handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (h *handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.Hello).Methods("GET")
	router.HandleFunc("/upload", h.uploadObjectHandler.UploadObject).Methods("POST")
}