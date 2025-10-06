package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormLogger() logger.Interface {
	// Custom logger output
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // log ke stdout
		logger.Config{
			SlowThreshold:             100 * time.Millisecond, // threshold 100ms untuk slow query
			LogLevel:                  logger.Info,            // tampilkan semua query (Info, Warn, Error)
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	return newLogger
}

func NewDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	log.Println("Connected to database")

	return db
}
