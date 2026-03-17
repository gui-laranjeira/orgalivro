package db

import (
	"gorm.io/gorm"

	"orgalivro/backend/internal/model"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Profile{},
		&model.Author{},
		&model.Genre{},
		&model.Book{},
		&model.LibraryEntry{},
	)
}
