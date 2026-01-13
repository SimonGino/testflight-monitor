package repository

import (
	"os"
	"path/filepath"

	"tf-monitor/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// Auto migrate tables
	return DB.AutoMigrate(
		&model.Monitor{},
		&model.TelegramConfig{},
		&model.SystemConfig{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
