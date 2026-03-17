package db

import (
	"gorm.io/gorm"

	"orgalivro/backend/internal/model"
)

func Migrate(db *gorm.DB) error {
	// Disable FK enforcement during migration — GORM's SQLite migrator may
	// recreate tables (drop + recreate), which fails when referencing tables
	// have FK constraints. Re-enable afterwards.
	db.Exec("PRAGMA foreign_keys=OFF")
	err := db.AutoMigrate(
		&model.Profile{},
		&model.Author{},
		&model.Genre{},
		&model.Book{},
		&model.LibraryEntry{},
	)
	db.Exec("PRAGMA foreign_keys=ON")
	return err
}
