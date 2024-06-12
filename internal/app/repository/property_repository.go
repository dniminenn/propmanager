package repository

import (
	"propmanager/internal/app/model"

	"gorm.io/gorm"
)

type PropertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

func (r *PropertyRepository) GetAllProperties() ([]model.Property, error) {
	var properties []model.Property
	err := r.db.Model(&model.Property{}).Preload("Images").Find(&properties).Error
	return properties, err
}

func (r *PropertyRepository) GetProperty(id uint) (model.Property, error) {
	var property model.Property
	err := r.db.Model(&model.Property{}).Preload("Images").First(&property, id).Error
	return property, err
}

func (r *PropertyRepository) CreateProperty(property *model.Property) error {
	return r.db.Create(property).Error
}

func (r *PropertyRepository) UpdateProperty(property *model.Property) error {
	return r.db.Save(property).Error
}

func (r *PropertyRepository) DeleteProperty(id uint) error {
	return r.db.Delete(&model.Property{}, id).Error
}

func (r *PropertyRepository) DeleteImage(propertyID uint, imageID uint) error {
	return r.db.Where("property_id = ? AND id = ?", propertyID, imageID).Delete(&model.Image{}).Error
}
