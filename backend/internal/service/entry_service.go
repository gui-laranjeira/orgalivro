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
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return s.repo.GetByProfileAndBook(profileID, req.BookID)
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
	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	// Reload to get full book preloads
	e.Book = model.Book{}
	return s.repo.GetByProfileAndBook(profileID, bookID)
}

func (s *EntryService) Remove(profileID, bookID uint) error {
	return s.repo.Delete(profileID, bookID)
}
