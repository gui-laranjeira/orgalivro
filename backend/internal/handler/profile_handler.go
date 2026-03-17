package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"orgalivro/backend/internal/dto"
	"orgalivro/backend/internal/model"
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
	resp := make([]dto.ProfileResponse, len(profiles))
	for i, p := range profiles {
		resp[i] = toProfileResponse(p)
	}
	c.JSON(http.StatusOK, resp)
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
	c.JSON(http.StatusCreated, toProfileResponse(*p))
}

func (h *ProfileHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func toProfileResponse(p model.Profile) dto.ProfileResponse {
	createdAt := p.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	return dto.ProfileResponse{
		ID:        p.ID,
		Name:      p.Name,
		AvatarURL: p.AvatarURL,
		CreatedAt: createdAt.Format(time.RFC3339),
	}
}
