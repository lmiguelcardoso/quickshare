package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	httphandler "quickshare/adapter/http"
	"quickshare/adapter/repository"
	"quickshare/core/service"
	"quickshare/internal/config"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.NewConfig()

	log.Println("Starting QuickShare Backend...")
	log.Printf("Environment detected: DB_HOST=%s", cfg.DBConfig.DBHost)

	var db *sql.DB
	if cfg.ServerConfig.Port == "3000" {
		db = connectWithRetry(cfg, 5, 3*time.Second)
		defer db.Close()
	}

	// Initialize repositories
	postgresRepo := repository.NewPostgreSQLRepository(db)
	
	s3BlobStorage, err := repository.NewS3BlobStorage(
		cfg.S3Config.Region,
		cfg.S3Config.Bucket,
		cfg.S3Config.AccessKeyID,
		cfg.S3Config.SecretAccessKey,
	)
	if err != nil {
		log.Fatal("Failed to initialize S3:", err)
	}

	// Initialize service
	uploadObjectService := service.NewUploadObjectService(postgresRepo, s3BlobStorage)

	// Initialize handlers
	uploadObjectHandler := httphandler.NewUploadObjectHandler(uploadObjectService)
	handler := httphandler.NewHandler(uploadObjectHandler)

	// Setup routes
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	log.Printf("Server is running on port %s", cfg.ServerConfig.Port)

	if err := http.ListenAndServe(":" + cfg.ServerConfig.Port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func connectWithRetry(cfg *config.Config, maxRetries int, delay time.Duration) *sql.DB {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBConfig.DBHost,
		cfg.DBConfig.DBPort,
		cfg.DBConfig.DBUser,
		cfg.DBConfig.DBPassword,
		cfg.DBConfig.DBName,
	)

	log.Printf("Attempting to connect to database at %s:%s...", cfg.DBConfig.DBHost, cfg.DBConfig.DBPort)

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Printf("Retry %d/%d: Failed to open database connection: %v", i+1, maxRetries, err)
			time.Sleep(delay)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database!")
			return db
		}

		log.Printf("Retry %d/%d: Failed to ping database: %v", i+1, maxRetries, err)
		db.Close()
		time.Sleep(delay)
	}

	log.Fatal("Failed to connect to database after retries:", err)
	return nil
}