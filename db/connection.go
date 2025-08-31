package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DB_NAME = "test_db_30Aug2025.sqlite3"

func SQLiteConnect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
