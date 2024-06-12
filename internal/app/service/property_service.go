package service

import (
	"propmanager/internal/app/model"
	"propmanager/internal/app/repository"
)

type PropertyService struct {
	repo *repository.PropertyRepository
}

func NewPropertyService(repo *repository.PropertyRepository) *PropertyService {
	return &PropertyService{repo: repo}
}

func (s *PropertyService) GetAllProperties() ([]model.Property, error) {
	return s.repo.GetAllProperties()
}

func (s *PropertyService) GetProperty(id uint) (model.Property, error) {
	return s.repo.GetProperty(id)
}

func (s *PropertyService) CreateProperty(property *model.Property) error {
	return s.repo.CreateProperty(property)
}

func (s *PropertyService) UpdateProperty(property *model.Property) error {
	return s.repo.UpdateProperty(property)
}

func (s *PropertyService) DeleteProperty(id uint) error {
	return s.repo.DeleteProperty(id)
}

func (s *PropertyService) DeleteImage(propertyID uint, imageID uint) error {
	return s.repo.DeleteImage(propertyID, imageID)
}
