package service

import (
	"errors"

	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
	"orgalivro/backend/internal/repository"
)

var ErrBookHasEntries = errors.New("book has library entries and cannot be deleted")

type BookService struct{ repo *repository.BookRepo }

func NewBookService(r *repository.BookRepo) *BookService { return &BookService{r} }

func (s *BookService) List(q dto.BookListQuery) ([]model.Book, int64, error) {
	return s.repo.List(q)
}

func (s *BookService) GetByID(id uint) (*model.Book, error) { return s.repo.GetByID(id) }

func (s *BookService) TransferOwner(bookID, profileID uint) error {
	return s.repo.SetOwner(bookID, profileID)
}

func (s *BookService) Create(req dto.CreateBookRequest) (*model.Book, error) {
	var ownerID *uint
	if req.OwnerProfileID != 0 {
		id := req.OwnerProfileID
		ownerID = &id
	}
	book := &model.Book{
		Title:          req.Title,
		ISBN13:         req.ISBN13,
		CoverURL:       req.CoverURL,
		Description:    req.Description,
		Year:           req.Year,
		Language:       req.Language,
		OwnerProfileID: ownerID,
	}
	for _, name := range req.Authors {
		a, err := s.repo.FindOrCreateAuthor(name)
		if err != nil {
			return nil, err
		}
		book.Authors = append(book.Authors, *a)
	}
	for _, name := range req.Genres {
		g, err := s.repo.FindOrCreateGenre(name)
		if err != nil {
			return nil, err
		}
		book.Genres = append(book.Genres, *g)
	}
	return book, s.repo.Create(book)
}

func (s *BookService) Update(id uint, req dto.UpdateBookRequest) (*model.Book, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.ISBN13 != "" {
		book.ISBN13 = req.ISBN13
	}
	if req.CoverURL != "" {
		book.CoverURL = req.CoverURL
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Year != 0 {
		book.Year = req.Year
	}
	if req.Language != "" {
		book.Language = req.Language
	}
	if req.Authors != nil {
		book.Authors = nil
		for _, name := range req.Authors {
			a, err := s.repo.FindOrCreateAuthor(name)
			if err != nil {
				return nil, err
			}
			book.Authors = append(book.Authors, *a)
		}
	}
	if req.Genres != nil {
		book.Genres = nil
		for _, name := range req.Genres {
			g, err := s.repo.FindOrCreateGenre(name)
			if err != nil {
				return nil, err
			}
			book.Genres = append(book.Genres, *g)
		}
	}
	return book, s.repo.Update(book)
}

func (s *BookService) Delete(id uint) error {
	has, err := s.repo.HasEntries(id)
	if err != nil {
		return err
	}
	if has {
		return ErrBookHasEntries
	}
	return s.repo.Delete(id)
}

func (s *BookService) AllAuthors() ([]model.Author, error) { return s.repo.AllAuthors() }
func (s *BookService) AllGenres() ([]model.Genre, error)   { return s.repo.AllGenres() }
