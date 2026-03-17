package model

import "time"

type Profile struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	AvatarURL string
	CreatedAt time.Time
}
