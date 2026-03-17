package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
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
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 20
	}
	books, total, err := h.svc.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedBooks{
		Data:  toBookResponses(books),
		Total: total,
		Page:  q.Page,
		Limit: q.Limit,
	})
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

func (h *BookHandler) TransferOwner(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	var req dto.TransferOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.TransferOwner(id, req.ProfileID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	book, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toBookResponse(*book))
}

func (h *BookHandler) Authors(c *gin.Context) {
	authors, err := h.svc.AllAuthors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]dto.AuthorResponse, len(authors))
	for i, a := range authors {
		resp[i] = dto.AuthorResponse{ID: a.ID, Name: a.Name}
	}
	c.JSON(http.StatusOK, resp)
}

func (h *BookHandler) Genres(c *gin.Context) {
	genres, err := h.svc.AllGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]dto.GenreResponse, len(genres))
	for i, g := range genres {
		resp[i] = dto.GenreResponse{ID: g.ID, Name: g.Name}
	}
	c.JSON(http.StatusOK, resp)
}

// parseID is shared across all handlers in this package.
func parseID(c *gin.Context, key string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
	}
	return uint(id), err
}

func toBookResponse(b model.Book) dto.BookResponse {
	authors := make([]dto.AuthorResponse, len(b.Authors))
	for i, a := range b.Authors {
		authors[i] = dto.AuthorResponse{ID: a.ID, Name: a.Name}
	}
	genres := make([]dto.GenreResponse, len(b.Genres))
	for i, g := range b.Genres {
		genres[i] = dto.GenreResponse{ID: g.ID, Name: g.Name}
	}
	createdAt := b.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	updatedAt := b.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}
	readers := make([]dto.BookReader, len(b.LibraryEntries))
	for i, e := range b.LibraryEntries {
		readers[i] = dto.BookReader{
			ProfileID:   e.ProfileID,
			ProfileName: e.Profile.Name,
			Status:      e.Status,
		}
	}
	return dto.BookResponse{
		ID:               b.ID,
		Title:            b.Title,
		ISBN13:           b.ISBN13,
		CoverURL:         b.CoverURL,
		Description:      b.Description,
		Year:             b.Year,
		Language:         b.Language,
		OwnerProfileID:   b.OwnerProfileID,
		OwnerProfileName: b.OwnerProfile.Name,
		Authors:          authors,
		Genres:           genres,
		Readers:          readers,
		CreatedAt:        createdAt.Format(time.RFC3339),
		UpdatedAt:        updatedAt.Format(time.RFC3339),
	}
}

func toBookResponses(books []model.Book) []dto.BookResponse {
	resp := make([]dto.BookResponse, len(books))
	for i, b := range books {
		resp[i] = toBookResponse(b)
	}
	return resp
}
