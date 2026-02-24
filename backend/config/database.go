package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database はDB接続を保持する構造体
type Database struct {
	DB *gorm.DB
}

// NewDatabase はDB接続を確立し、Database構造体を返す
func NewDatabase() (*Database, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "caltrack")
	password := getEnv("DB_PASSWORD", "caltrack")
	dbname := getEnv("DB_NAME", "caltrack")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return &Database{DB: db}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
