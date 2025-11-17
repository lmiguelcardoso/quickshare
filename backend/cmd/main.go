package main

import (
	"log"
	"net/http"
	httphandler "quickshare/adapter/http"

	"github.com/gorilla/mux"
)

const port = ":8080"

func main() {
	handler := httphandler.NewHandler()
	router := mux.NewRouter()
	
	handler.RegisterRoutes(router)
	log.Println("Server is running on port", port)

	err := http.ListenAndServe(port, router); if err != nil {
		log.Fatal(err)
	}
}