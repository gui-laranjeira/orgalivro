package repository

import (
	"gorm.io/gorm"

	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
)

type BookRepo struct{ db *gorm.DB }

func NewBookRepo(db *gorm.DB) *BookRepo { return &BookRepo{db} }

func (r *BookRepo) List(q dto.BookListQuery) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	tx := r.db.Model(&model.Book{})

	if q.Q != "" {
		tx = tx.Where("title LIKE ? OR description LIKE ?", "%"+q.Q+"%", "%"+q.Q+"%")
	}
	if q.Author != "" {
		tx = tx.Joins("JOIN book_authors ba ON ba.book_id = books.id").
			Joins("JOIN authors a ON a.id = ba.author_id").
			Where("a.name LIKE ?", "%"+q.Author+"%")
	}
	if q.Genre != "" {
		tx = tx.Joins("JOIN book_genres bg ON bg.book_id = books.id").
			Joins("JOIN genres g ON g.id = bg.genre_id").
			Where("g.name LIKE ?", "%"+q.Genre+"%")
	}
	if q.Year != 0 {
		tx = tx.Where("year = ?", q.Year)
	}
	if q.Language != "" {
		tx = tx.Where("language = ?", q.Language)
	}
	if q.OwnerProfileID != 0 {
		tx = tx.Where("owner_profile_id = ?", q.OwnerProfileID)
	}

	tx.Count(&total)

	offset := (q.Page - 1) * q.Limit
	err := tx.Preload("OwnerProfile").Preload("Authors").Preload("Genres").Preload("LibraryEntries.Profile").
		Offset(offset).Limit(q.Limit).Order("title").Find(&books).Error
	return books, total, err
}

func (r *BookRepo) GetByID(id uint) (*model.Book, error) {
	var book model.Book
	err := r.db.Preload("OwnerProfile").Preload("Authors").Preload("Genres").Preload("LibraryEntries.Profile").First(&book, id).Error
	return &book, err
}

func (r *BookRepo) SetOwner(bookID, profileID uint) error {
	return r.db.Model(&model.Book{}).Where("id = ?", bookID).Update("owner_profile_id", profileID).Error
}

func (r *BookRepo) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepo) Update(book *model.Book) error {
	if err := r.db.Model(book).Association("Authors").Replace(book.Authors); err != nil {
		return err
	}
	if err := r.db.Model(book).Association("Genres").Replace(book.Genres); err != nil {
		return err
	}
	return r.db.Save(book).Error
}

func (r *BookRepo) Delete(id uint) error {
	return r.db.Delete(&model.Book{}, id).Error
}

func (r *BookRepo) HasEntries(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.LibraryEntry{}).Where("book_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *BookRepo) FindOrCreateAuthor(name string) (*model.Author, error) {
	var a model.Author
	err := r.db.Where(model.Author{Name: name}).FirstOrCreate(&a).Error
	return &a, err
}

func (r *BookRepo) FindOrCreateGenre(name string) (*model.Genre, error) {
	var g model.Genre
	err := r.db.Where(model.Genre{Name: name}).FirstOrCreate(&g).Error
	return &g, err
}

func (r *BookRepo) AllAuthors() ([]model.Author, error) {
	var authors []model.Author
	return authors, r.db.Order("name").Find(&authors).Error
}

func (r *BookRepo) AllGenres() ([]model.Genre, error) {
	var genres []model.Genre
	return genres, r.db.Order("name").Find(&genres).Error
}
