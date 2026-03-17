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
