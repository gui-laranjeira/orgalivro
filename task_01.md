# Task 01 — Backend Implementation

Go + Gin + GORM + SQLite backend for Orgalivro.

---

## Prerequisites

- Go 1.22+
- `air` for hot reload: `go install github.com/air-verse/air@latest`

---

## Step 1 — Initialize Go module and install dependencies

```bash
mkdir -p backend && cd backend
go mod init orgalivro/backend
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/glebarez/sqlite        # pure-Go SQLite driver (CGO-free)
go get github.com/glebarez/gorm-sqlite   # GORM dialect wrapping glebarez/sqlite
```

> Use `github.com/glebarez/gorm-sqlite` as the GORM dialect. It is a drop-in replacement for `gorm.io/driver/sqlite` with no CGO requirement.

---

## Step 2 — Create directory tree

```bash
mkdir -p cmd/server
mkdir -p internal/{config,db,model,repository,service,handler,router,dto}
```

---

## Step 3 — `internal/config/config.go`

```go
package config

import (
	"os"
	"strings"
)

type Config struct {
	DBPath         string
	Port           string
	AllowedOrigins []string
}

func Load() Config {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		origins = "http://localhost:5173"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./orgalivro.db"
	}
	return Config{
		DBPath:         dbPath,
		Port:           port,
		AllowedOrigins: strings.Split(origins, ","),
	}
}
```

---

## Step 4 — GORM Models

### `internal/model/profile.go`

```go
package model

import "time"

type Profile struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	AvatarURL string
	CreatedAt time.Time
}
```

### `internal/model/book.go`

```go
package model

import "time"

type Author struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"uniqueIndex;not null"`
	Books []Book `gorm:"many2many:book_authors;"`
}

type Genre struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"uniqueIndex;not null"`
	Books []Book `gorm:"many2many:book_genres;"`
}

type Book struct {
	ID          uint      `gorm:"primarykey"`
	Title       string    `gorm:"not null"`
	ISBN13      string    `gorm:"uniqueIndex"`
	CoverURL    string
	Description string
	Year        int
	Language    string
	Authors     []Author  `gorm:"many2many:book_authors;"`
	Genres      []Genre   `gorm:"many2many:book_genres;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

### `internal/model/entry.go`

```go
package model

import "time"

type LibraryEntry struct {
	ID        uint    `gorm:"primarykey"`
	ProfileID uint    `gorm:"not null;uniqueIndex:idx_profile_book"`
	BookID    uint    `gorm:"not null;uniqueIndex:idx_profile_book"`
	Status    string  `gorm:"not null;default:'want_to_read'"`
	Rating    *int
	Notes     string
	AddedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Profile   Profile
	Book      Book
}
```

---

## Step 5 — Database Setup

### `internal/db/db.go`

```go
package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	glebarez "github.com/glebarez/gorm-sqlite"
)

func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(glebarez.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)

	// Enable WAL mode and foreign keys
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	return db, nil
}
```

### `internal/db/migrate.go`

```go
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
```

---

## Step 6 — DTOs

### `internal/dto/profile_dto.go`

```go
package dto

type CreateProfileRequest struct {
	Name      string `json:"name" binding:"required"`
	AvatarURL string `json:"avatar_url"`
}

type ProfileResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
}
```

### `internal/dto/book_dto.go`

```go
package dto

type CreateBookRequest struct {
	Title       string   `json:"title" binding:"required"`
	ISBN13      string   `json:"isbn13"`
	CoverURL    string   `json:"cover_url"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Language    string   `json:"language"`
	Authors     []string `json:"authors"`
	Genres      []string `json:"genres"`
}

type UpdateBookRequest struct {
	Title       string   `json:"title"`
	ISBN13      string   `json:"isbn13"`
	CoverURL    string   `json:"cover_url"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Language    string   `json:"language"`
	Authors     []string `json:"authors"`
	Genres      []string `json:"genres"`
}

type BookListQuery struct {
	Q        string `form:"q"`
	Author   string `form:"author"`
	Genre    string `form:"genre"`
	Year     int    `form:"year"`
	Language string `form:"language"`
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=20"`
}

type AuthorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookResponse struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	ISBN13      string           `json:"isbn13"`
	CoverURL    string           `json:"cover_url"`
	Description string           `json:"description"`
	Year        int              `json:"year"`
	Language    string           `json:"language"`
	Authors     []AuthorResponse `json:"authors"`
	Genres      []GenreResponse  `json:"genres"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

type PaginatedBooks struct {
	Data  []BookResponse `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
```

### `internal/dto/entry_dto.go`

```go
package dto

type CreateEntryRequest struct {
	BookID uint   `json:"book_id" binding:"required"`
	Status string `json:"status"`
	Rating *int   `json:"rating"`
	Notes  string `json:"notes"`
}

type UpdateEntryRequest struct {
	Status string `json:"status"`
	Rating *int   `json:"rating"`
	Notes  string `json:"notes"`
}

type EntryListQuery struct {
	Status string `form:"status"`
	Q      string `form:"q"`
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=20"`
}

type EntryResponse struct {
	ID        uint         `json:"id"`
	ProfileID uint         `json:"profile_id"`
	Book      BookResponse `json:"book"`
	Status    string       `json:"status"`
	Rating    *int         `json:"rating"`
	Notes     string       `json:"notes"`
	AddedAt   string       `json:"added_at"`
	UpdatedAt string       `json:"updated_at"`
}

type PaginatedEntries struct {
	Data  []EntryResponse `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}
```

---

## Step 7 — Repositories

### `internal/repository/profile_repo.go`

```go
package repository

import (
	"gorm.io/gorm"
	"orgalivro/backend/internal/model"
)

type ProfileRepo struct{ db *gorm.DB }

func NewProfileRepo(db *gorm.DB) *ProfileRepo { return &ProfileRepo{db} }

func (r *ProfileRepo) List() ([]model.Profile, error) {
	var profiles []model.Profile
	return profiles, r.db.Order("name").Find(&profiles).Error
}

func (r *ProfileRepo) Create(p *model.Profile) error {
	return r.db.Create(p).Error
}

func (r *ProfileRepo) Delete(id uint) error {
	return r.db.Delete(&model.Profile{}, id).Error
}

func (r *ProfileRepo) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Profile{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
```

### `internal/repository/book_repo.go`

```go
package repository

import (
	"gorm.io/gorm"
	"orgalivro/backend/internal/model"
	"orgalivro/backend/internal/dto"
)

type BookRepo struct{ db *gorm.DB }

func NewBookRepo(db *gorm.DB) *BookRepo { return &BookRepo{db} }

func (r *BookRepo) List(q dto.BookListQuery) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	tx := r.db.Model(&model.Book{}).Preload("Authors").Preload("Genres")

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

	tx.Count(&total)
	offset := (q.Page - 1) * q.Limit
	err := tx.Offset(offset).Limit(q.Limit).Order("title").Find(&books).Error
	return books, total, err
}

func (r *BookRepo) GetByID(id uint) (*model.Book, error) {
	var book model.Book
	err := r.db.Preload("Authors").Preload("Genres").First(&book, id).Error
	return &book, err
}

func (r *BookRepo) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepo) Update(book *model.Book) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(book).Error
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
```

> Note: `model.LibraryEntry` is referenced here for the `HasEntries` check — import accordingly.

### `internal/repository/entry_repo.go`

```go
package repository

import (
	"gorm.io/gorm"
	"orgalivro/backend/internal/model"
	"orgalivro/backend/internal/dto"
)

type EntryRepo struct{ db *gorm.DB }

func NewEntryRepo(db *gorm.DB) *EntryRepo { return &EntryRepo{db} }

func (r *EntryRepo) List(profileID uint, q dto.EntryListQuery) ([]model.LibraryEntry, int64, error) {
	var entries []model.LibraryEntry
	var total int64

	tx := r.db.Model(&model.LibraryEntry{}).
		Where("profile_id = ?", profileID).
		Preload("Book.Authors").
		Preload("Book.Genres")

	if q.Status != "" {
		tx = tx.Where("status = ?", q.Status)
	}
	if q.Q != "" {
		tx = tx.Joins("JOIN books b ON b.id = library_entries.book_id").
			Where("b.title LIKE ?", "%"+q.Q+"%")
	}

	tx.Count(&total)
	offset := (q.Page - 1) * q.Limit
	err := tx.Offset(offset).Limit(q.Limit).Order("updated_at DESC").Find(&entries).Error
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
	return r.db.Save(e).Error
}

func (r *EntryRepo) Delete(profileID, bookID uint) error {
	return r.db.
		Where("profile_id = ? AND book_id = ?", profileID, bookID).
		Delete(&model.LibraryEntry{}).Error
}
```

---

## Step 8 — Services

### `internal/service/profile_service.go`

```go
package service

import (
	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
	"orgalivro/backend/internal/repository"
)

type ProfileService struct{ repo *repository.ProfileRepo }

func NewProfileService(r *repository.ProfileRepo) *ProfileService { return &ProfileService{r} }

func (s *ProfileService) List() ([]model.Profile, error) { return s.repo.List() }

func (s *ProfileService) Create(req dto.CreateProfileRequest) (*model.Profile, error) {
	p := &model.Profile{Name: req.Name, AvatarURL: req.AvatarURL}
	return p, s.repo.Create(p)
}

func (s *ProfileService) Delete(id uint) error { return s.repo.Delete(id) }
```

### `internal/service/book_service.go`

```go
package service

import (
	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
	"orgalivro/backend/internal/repository"
)

type BookService struct{ repo *repository.BookRepo }

func NewBookService(r *repository.BookRepo) *BookService { return &BookService{r} }

func (s *BookService) List(q dto.BookListQuery) ([]model.Book, int64, error) {
	return s.repo.List(q)
}

func (s *BookService) GetByID(id uint) (*model.Book, error) { return s.repo.GetByID(id) }

func (s *BookService) Create(req dto.CreateBookRequest) (*model.Book, error) {
	book := &model.Book{
		Title:       req.Title,
		ISBN13:      req.ISBN13,
		CoverURL:    req.CoverURL,
		Description: req.Description,
		Year:        req.Year,
		Language:    req.Language,
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

// sentinel error
var ErrBookHasEntries = fmt.Errorf("book has library entries and cannot be deleted")
```

> Add `"fmt"` to imports.

### `internal/service/entry_service.go`

```go
package service

import (
	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
	"orgalivro/backend/internal/repository"
)

type EntryService struct {
	repo        *repository.EntryRepo
	profileRepo *repository.ProfileRepo
}

func NewEntryService(r *repository.EntryRepo, pr *repository.ProfileRepo) *EntryService {
	return &EntryService{r, pr}
}

func (s *EntryService) List(profileID uint, q dto.EntryListQuery) ([]model.LibraryEntry, int64, error) {
	return s.repo.List(profileID, q)
}

func (s *EntryService) Add(profileID uint, req dto.CreateEntryRequest) (*model.LibraryEntry, error) {
	status := req.Status
	if status == "" {
		status = "want_to_read"
	}
	e := &model.LibraryEntry{
		ProfileID: profileID,
		BookID:    req.BookID,
		Status:    status,
		Rating:    req.Rating,
		Notes:     req.Notes,
	}
	return e, s.repo.Create(e)
}

func (s *EntryService) Update(profileID, bookID uint, req dto.UpdateEntryRequest) (*model.LibraryEntry, error) {
	e, err := s.repo.GetByProfileAndBook(profileID, bookID)
	if err != nil {
		return nil, err
	}
	if req.Status != "" {
		e.Status = req.Status
	}
	if req.Rating != nil {
		e.Rating = req.Rating
	}
	if req.Notes != "" {
		e.Notes = req.Notes
	}
	return e, s.repo.Update(e)
}

func (s *EntryService) Remove(profileID, bookID uint) error {
	return s.repo.Delete(profileID, bookID)
}
```

### `internal/service/isbn_service.go`

```go
package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ISBNResult struct {
	Title       string   `json:"title"`
	ISBN13      string   `json:"isbn13"`
	CoverURL    string   `json:"cover_url"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Language    string   `json:"language"`
	Authors     []string `json:"authors"`
	Genres      []string `json:"genres"`
}

type ISBNService struct{ client *http.Client }

func NewISBNService() *ISBNService {
	return &ISBNService{client: &http.Client{}}
}

func (s *ISBNService) Lookup(isbn string) (*ISBNResult, error) {
	url := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", isbn)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	key := "ISBN:" + isbn
	entry, ok := raw[key]
	if !ok {
		return nil, fmt.Errorf("ISBN not found")
	}

	data := entry.(map[string]any)
	result := &ISBNResult{ISBN13: isbn}

	if v, ok := data["title"].(string); ok {
		result.Title = v
	}

	// Cover
	if covers, ok := data["cover"].(map[string]any); ok {
		if large, ok := covers["large"].(string); ok {
			result.CoverURL = large
		} else if medium, ok := covers["medium"].(string); ok {
			result.CoverURL = medium
		}
	}

	// Publish year
	if py, ok := data["publish_date"].(string); ok {
		// Extract 4-digit year
		for i := 0; i+4 <= len(py); i++ {
			year := 0
			fmt.Sscanf(py[i:i+4], "%d", &year)
			if year > 1000 && year < 2100 {
				result.Year = year
				break
			}
		}
	}

	// Authors
	if authors, ok := data["authors"].([]any); ok {
		for _, a := range authors {
			if am, ok := a.(map[string]any); ok {
				if name, ok := am["name"].(string); ok {
					result.Authors = append(result.Authors, name)
				}
			}
		}
	}

	// Description
	if notes, ok := data["notes"].(map[string]any); ok {
		if v, ok := notes["value"].(string); ok {
			result.Description = v
		}
	} else if notes, ok := data["notes"].(string); ok {
		result.Description = notes
	}

	// Language from subjects or just leave empty
	if subjects, ok := data["subjects"].([]any); ok {
		genres := []string{}
		for _, s := range subjects {
			if sm, ok := s.(map[string]any); ok {
				if name, ok := sm["name"].(string); ok {
					genres = append(genres, strings.Title(name))
				}
			}
		}
		result.Genres = genres
	}

	return result, nil
}
```

---

## Step 9 — Handlers

### `internal/handler/profile_handler.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/service"
)

type ProfileHandler struct{ svc *service.ProfileService }

func NewProfileHandler(svc *service.ProfileService) *ProfileHandler { return &ProfileHandler{svc} }

func (h *ProfileHandler) List(c *gin.Context) {
	profiles, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *ProfileHandler) Create(c *gin.Context) {
	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.svc.Create(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProfileHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
```

### `internal/handler/book_handler.go`

```go
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/service"
)

type BookHandler struct{ svc *service.BookService }

func NewBookHandler(svc *service.BookService) *BookHandler { return &BookHandler{svc} }

func (h *BookHandler) List(c *gin.Context) {
	var q dto.BookListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books, total, err := h.svc.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedBooks{Data: toBookResponses(books), Total: total, Page: q.Page, Limit: q.Limit})
}

func (h *BookHandler) Get(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	book, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, toBookResponse(*book))
}

func (h *BookHandler) Create(c *gin.Context) {
	var req dto.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book, err := h.svc.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toBookResponse(*book))
}

func (h *BookHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	var req dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book, err := h.svc.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toBookResponse(*book))
}

func (h *BookHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrBookHasEntries) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BookHandler) Authors(c *gin.Context) {
	authors, err := h.svc.AllAuthors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, authors)
}

func (h *BookHandler) Genres(c *gin.Context) {
	genres, err := h.svc.AllGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, genres)
}

// helpers
func parseID(c *gin.Context, key string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
	}
	return uint(id), err
}

func toBookResponse(b model.Book) dto.BookResponse { /* map fields */ }
func toBookResponses(books []model.Book) []dto.BookResponse { /* map slice */ }
```

> Implement `toBookResponse` and `toBookResponses` by mapping model fields to DTO fields. Import `orgalivro/backend/internal/model`.

### `internal/handler/entry_handler.go`

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/service"
)

type EntryHandler struct{ svc *service.EntryService }

func NewEntryHandler(svc *service.EntryService) *EntryHandler { return &EntryHandler{svc} }

func (h *EntryHandler) List(c *gin.Context) {
	pid, err := parseID(c, "pid")
	if err != nil {
		return
	}
	var q dto.EntryListQuery
	c.ShouldBindQuery(&q)
	entries, total, err := h.svc.List(pid, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedEntries{Data: toEntryResponses(entries), Total: total, Page: q.Page, Limit: q.Limit})
}

func (h *EntryHandler) Add(c *gin.Context) {
	pid, err := parseID(c, "pid")
	if err != nil {
		return
	}
	var req dto.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entry, err := h.svc.Add(pid, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *EntryHandler) Update(c *gin.Context) {
	pid, err := parseID(c, "pid")
	if err != nil {
		return
	}
	bookID, err := parseID(c, "book_id")
	if err != nil {
		return
	}
	var req dto.UpdateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entry, err := h.svc.Update(pid, bookID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) Remove(c *gin.Context) {
	pid, err := parseID(c, "pid")
	if err != nil {
		return
	}
	bookID, err := parseID(c, "book_id")
	if err != nil {
		return
	}
	if err := h.svc.Remove(pid, bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func toEntryResponse(e model.LibraryEntry) dto.EntryResponse { /* map fields */ }
func toEntryResponses(entries []model.LibraryEntry) []dto.EntryResponse { /* map slice */ }
```

### `internal/handler/isbn_handler.go`

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"orgalivro/backend/internal/service"
)

type ISBNHandler struct{ svc *service.ISBNService }

func NewISBNHandler(svc *service.ISBNService) *ISBNHandler { return &ISBNHandler{svc} }

func (h *ISBNHandler) Lookup(c *gin.Context) {
	isbn := c.Param("isbn")
	result, err := h.svc.Lookup(isbn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
```

---

## Step 10 — Router

### `internal/router/router.go`

```go
package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"orgalivro/backend/internal/config"
	"orgalivro/backend/internal/handler"
)

func New(
	cfg config.Config,
	ph *handler.ProfileHandler,
	bh *handler.BookHandler,
	eh *handler.EntryHandler,
	ih *handler.ISBNHandler,
) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.AllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(corsConfig))

	v1 := r.Group("/api/v1")
	{
		// Profiles
		v1.GET("/profiles", ph.List)
		v1.POST("/profiles", ph.Create)
		v1.DELETE("/profiles/:id", ph.Delete)

		// Books
		v1.GET("/books", bh.List)
		v1.POST("/books", bh.Create)
		v1.GET("/books/:id", bh.Get)
		v1.PUT("/books/:id", bh.Update)
		v1.DELETE("/books/:id", bh.Delete)
		v1.GET("/authors", bh.Authors)
		v1.GET("/genres", bh.Genres)

		// Library entries
		v1.GET("/profiles/:pid/library", eh.List)
		v1.POST("/profiles/:pid/library", eh.Add)
		v1.PUT("/profiles/:pid/library/:book_id", eh.Update)
		v1.DELETE("/profiles/:pid/library/:book_id", eh.Remove)

		// ISBN lookup
		v1.GET("/isbn/:isbn", ih.Lookup)
	}

	return r
}
```

---

## Step 11 — `cmd/server/main.go`

```go
package main

import (
	"log"

	"orgalivro/backend/internal/config"
	"orgalivro/backend/internal/db"
	"orgalivro/backend/internal/handler"
	"orgalivro/backend/internal/repository"
	"orgalivro/backend/internal/router"
	"orgalivro/backend/internal/service"
)

func main() {
	cfg := config.Load()

	gormDB, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.Migrate(gormDB); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Repositories
	profileRepo := repository.NewProfileRepo(gormDB)
	bookRepo := repository.NewBookRepo(gormDB)
	entryRepo := repository.NewEntryRepo(gormDB)

	// Services
	profileSvc := service.NewProfileService(profileRepo)
	bookSvc := service.NewBookService(bookRepo)
	entrySvc := service.NewEntryService(entryRepo, profileRepo)
	isbnSvc := service.NewISBNService()

	// Handlers
	ph := handler.NewProfileHandler(profileSvc)
	bh := handler.NewBookHandler(bookSvc)
	eh := handler.NewEntryHandler(entrySvc)
	ih := handler.NewISBNHandler(isbnSvc)

	r := router.New(cfg, ph, bh, eh, ih)
	log.Printf("Starting server on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
```

---

## Step 12 — `.air.toml` (hot reload)

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false
```

---

## Step 13 — Smoke Tests

```bash
cd backend
go build ./cmd/server && ./server &

# List profiles (empty)
curl -s http://localhost:8080/api/v1/profiles | jq .

# Create a profile
curl -s -X POST http://localhost:8080/api/v1/profiles \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ana"}' | jq .

# Create a book
curl -s -X POST http://localhost:8080/api/v1/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"Dune","authors":["Frank Herbert"],"genres":["Science Fiction"],"year":1965}' | jq .

# Add book to Ana's library (use the IDs returned above)
curl -s -X POST http://localhost:8080/api/v1/profiles/1/library \
  -H 'Content-Type: application/json' \
  -d '{"book_id":1,"status":"reading"}' | jq .

# ISBN lookup
curl -s http://localhost:8080/api/v1/isbn/9780261102354 | jq .

# List Ana's library
curl -s http://localhost:8080/api/v1/profiles/1/library | jq .
```

---

## Notes

- `toBookResponse` / `toEntryResponse` helpers should map GORM model time fields to ISO8601 strings using `.Format(time.RFC3339)`.
- The `HasEntries` check in `BookRepo` references `model.LibraryEntry` — ensure the import is present.
- For GORM many2many updates (replacing associations on Update), use `db.Session(&gorm.Session{FullSaveAssociations: true}).Save(book)` and clear the slice before re-appending.
- Run `go mod tidy` after all files are created.
