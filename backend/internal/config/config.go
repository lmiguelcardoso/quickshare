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

type Config struct {
	DBConfig *DBConfig
	ServerConfig *ServerConfig
}

func NewConfig() *Config {
	var dbConfig DBConfig
	var serverConfig ServerConfig

	dbConfig.DBHost = os.Getenv("DB_HOST")
	dbConfig.DBPort = os.Getenv("DB_PORT")
	dbConfig.DBUser = os.Getenv("DB_USER")
	dbConfig.DBPassword = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")

	serverConfig.Port = os.Getenv("PORT")

	return &Config{
		DBConfig: &dbConfig,
		ServerConfig: &serverConfig,
	}
}