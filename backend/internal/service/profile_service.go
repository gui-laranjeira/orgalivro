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
