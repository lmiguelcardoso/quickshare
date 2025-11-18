package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	httphandler "quickshare/adapter/http"
	"quickshare/internal/config"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


func main() {
	handler := httphandler.NewHandler()
	router := mux.NewRouter()
	cfg := config.NewConfig()
	handler.RegisterRoutes(router)
	log.Println("Server is running on port", cfg.ServerConfig.Port)

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBConfig.DBHost, cfg.DBConfig.DBPort, cfg.DBConfig.DBUser, cfg.DBConfig.DBPassword, cfg.DBConfig.DBName)
	
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Connected to database")

	err = http.ListenAndServe(cfg.ServerConfig.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}