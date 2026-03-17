package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
	"orgalivro/backend/internal/service"
)

type EntryHandler struct{ svc *service.EntryService }

func NewEntryHandler(svc *service.EntryService) *EntryHandler { return &EntryHandler{svc} }

func (h *EntryHandler) List(c *gin.Context) {
	pid, err := parseID(c, "id")
	if err != nil {
		return
	}
	var q dto.EntryListQuery
	_ = c.ShouldBindQuery(&q)
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 20
	}
	entries, total, err := h.svc.List(pid, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedEntries{
		Data:  toEntryResponses(entries),
		Total: total,
		Page:  q.Page,
		Limit: q.Limit,
	})
}

func (h *EntryHandler) Add(c *gin.Context) {
	pid, err := parseID(c, "id")
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
	c.JSON(http.StatusCreated, toEntryResponse(*entry))
}

func (h *EntryHandler) Update(c *gin.Context) {
	pid, err := parseID(c, "id")
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
	c.JSON(http.StatusOK, toEntryResponse(*entry))
}

func (h *EntryHandler) Remove(c *gin.Context) {
	pid, err := parseID(c, "id")
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

func toEntryResponse(e model.LibraryEntry) dto.EntryResponse {
	addedAt := e.AddedAt
	if addedAt.IsZero() {
		addedAt = time.Now()
	}
	updatedAt := e.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}
	return dto.EntryResponse{
		ID:        e.ID,
		ProfileID: e.ProfileID,
		Book:      toBookResponse(e.Book),
		Status:    e.Status,
		Rating:    e.Rating,
		Notes:     e.Notes,
		AddedAt:   addedAt.Format(time.RFC3339),
		UpdatedAt: updatedAt.Format(time.RFC3339),
	}
}

func toEntryResponses(entries []model.LibraryEntry) []dto.EntryResponse {
	resp := make([]dto.EntryResponse, len(entries))
	for i, e := range entries {
		resp[i] = toEntryResponse(e)
	}
	return resp
}
