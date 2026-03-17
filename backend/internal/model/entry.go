package model

import "time"

type LibraryEntry struct {
	ID        uint      `gorm:"primarykey"`
	ProfileID uint      `gorm:"not null;uniqueIndex:idx_profile_book"`
	BookID    uint      `gorm:"not null;uniqueIndex:idx_profile_book"`
	Status    string    `gorm:"not null;default:'want_to_read'"`
	Rating    *int
	Notes     string
	AddedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Profile   Profile
	Book      Book
}
