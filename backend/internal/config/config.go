package config

import "os"

type DBConfig struct {
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
}


type ServerConfig struct {
	Port string
}

type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

type Config struct {
	DBConfig     *DBConfig
	ServerConfig *ServerConfig
	S3Config     *S3Config
}

func NewConfig() *Config {
	var dbConfig DBConfig
	var serverConfig ServerConfig
	var s3Config S3Config

	dbConfig.DBHost = os.Getenv("DB_HOST")
	dbConfig.DBPort = os.Getenv("DB_PORT")
	dbConfig.DBUser = os.Getenv("DB_USER")
	dbConfig.DBPassword = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")

	serverConfig.Port = os.Getenv("PORT")

	s3Config.Region = getEnvOrDefault("AWS_REGION", "us-east-2")
	s3Config.Bucket = getEnvOrDefault("AWS_BUCKET_NAME", "quickshare-assets")
	s3Config.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	s3Config.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	return &Config{
		DBConfig:     &dbConfig,
		ServerConfig: &serverConfig,
		S3Config:     &s3Config,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}