package model

import "time"

type Author struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"uniqueIndex;not null"`
	Books []Book `gorm:"many2many:book_authors;"`
}

type Genre struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"uniqueIndex;not null"`
	Books []Book `gorm:"many2many:book_genres;"`
}

type Book struct {
	ID             uint           `gorm:"primarykey"`
	Title          string         `gorm:"not null"`
	ISBN13         string         `gorm:"uniqueIndex"`
	CoverURL       string
	Description    string
	Year           int
	Language       string
	OwnerProfileID *uint
	OwnerProfile   Profile
	Authors        []Author       `gorm:"many2many:book_authors;"`
	Genres         []Genre        `gorm:"many2many:book_genres;"`
	LibraryEntries []LibraryEntry
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
