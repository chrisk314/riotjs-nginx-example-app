package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"backend/models"
)

var db *gorm.DB

//PostgresConfig stores postgres config.
type PostgresConfig struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
	SSLMode  string
}

// GetRequiredEnvVar tries to get env var referenced by key. Panics if unsuccesful.
func GetRequiredEnvVar(key string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	panic(fmt.Sprintf("Missing required env var: %s", key))
}

// GetPostgresConfig populates postgres config from env.
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

// migrate applies database auto migrations with Gorm.
func migrate() *gorm.DB {
	db.AutoMigrate(&models.Book{})
	return db
}

// InitDB performs database initialisation actions for `db`.
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
	migrate() // apply database migrations
	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}
	return db
}

// GetDB returns an initialised database.
func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	if err := db.DB().Ping(); err != nil {
		panic(fmt.Sprintf("Database connection test failed, err: %s", err))
	}
	return db
}
