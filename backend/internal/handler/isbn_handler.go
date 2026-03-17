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
