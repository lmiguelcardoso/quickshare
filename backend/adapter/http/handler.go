package http

import (
	"net/http"

	"github.com/gorilla/mux"
)


type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (h *handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", h.Hello).Methods("GET")
}