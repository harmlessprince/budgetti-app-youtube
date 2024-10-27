package common

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewMysql() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	log.Default().Println("Database connection successful")
	return db, nil
}
