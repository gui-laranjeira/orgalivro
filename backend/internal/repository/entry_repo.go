package repository

import (
	"gorm.io/gorm"

	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
)

type EntryRepo struct{ db *gorm.DB }

func NewEntryRepo(db *gorm.DB) *EntryRepo { return &EntryRepo{db} }

func (r *EntryRepo) List(profileID uint, q dto.EntryListQuery) ([]model.LibraryEntry, int64, error) {
	var entries []model.LibraryEntry
	var total int64

	tx := r.db.Model(&model.LibraryEntry{}).
		Where("profile_id = ?", profileID)

	if q.Status != "" {
		tx = tx.Where("status = ?", q.Status)
	}
	if q.Q != "" {
		tx = tx.Joins("JOIN books b ON b.id = library_entries.book_id").
			Where("b.title LIKE ?", "%"+q.Q+"%")
	}

	tx.Count(&total)

	offset := (q.Page - 1) * q.Limit
	err := tx.Preload("Book.Authors").Preload("Book.Genres").
		Offset(offset).Limit(q.Limit).Order("library_entries.updated_at DESC").Find(&entries).Error
	return entries, total, err
}

func (r *EntryRepo) Create(e *model.LibraryEntry) error {
	return r.db.Create(e).Error
}

func (r *EntryRepo) GetByProfileAndBook(profileID, bookID uint) (*model.LibraryEntry, error) {
	var e model.LibraryEntry
	err := r.db.
		Preload("Book.Authors").
		Preload("Book.Genres").
		Where("profile_id = ? AND book_id = ?", profileID, bookID).
		First(&e).Error
	return &e, err
}

func (r *EntryRepo) Update(e *model.LibraryEntry) error {
	return r.db.Omit("Profile", "Book").Save(e).Error
}

func (r *EntryRepo) Delete(profileID, bookID uint) error {
	return r.db.
		Where("profile_id = ? AND book_id = ?", profileID, bookID).
		Delete(&model.LibraryEntry{}).Error
}
