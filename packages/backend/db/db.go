package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() (*gorm.DB, error) {
	if db == nil {
		var err error
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
