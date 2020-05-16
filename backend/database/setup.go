package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"backend/models"
)

var db *gorm.DB

type PostgresConfig struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
	SSLMode  string
}

func GetRequiredEnvVar(key string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	panic(fmt.Sprintf("Missing required env var: %s", key))
}

func GetPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     GetRequiredEnvVar("POSTGRES_HOST"),
		Port:     GetRequiredEnvVar("POSTGRES_PORT"),
		DB:       GetRequiredEnvVar("POSTGRES_DB"),
		User:     GetRequiredEnvVar("POSTGRES_USER"),
		Password: GetRequiredEnvVar("POSTGRES_PASSWORD"),
		SSLMode:  GetRequiredEnvVar("POSTGRES_SSLMODE"),
	}
}

func migrate() *gorm.DB {
	db.AutoMigrate(&models.Book{})
	return db
}

func InitDB() *gorm.DB {
	dbConf := GetPostgresConfig()
	dbSpec := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		dbConf.Host, dbConf.Port, dbConf.DB, dbConf.User, dbConf.Password, dbConf.SSLMode,
	)
	conn, err := gorm.Open("postgres", dbSpec)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database with config: %s, err: %s", dbSpec, err))
	}
	db = conn
	migrate()
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}
