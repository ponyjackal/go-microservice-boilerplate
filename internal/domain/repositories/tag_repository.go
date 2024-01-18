package repositories

import (
	"strings"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/adapters/database"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/models"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/logger"
)

type TagRepository struct{}

func (r *TagRepository) Save(tag *models.Tag) error {
	err := database.DB.Create(tag).Error
	if err != nil {
		logger.Errorf("failed to save data: %v", err)
	}
	return err
}

func (r *TagRepository) GetTags(name string) ([]models.Tag, error) {
	var tags []models.Tag
	var err error

	db := database.DB
	if name != "" {
		nameQuery := "%" + strings.ToLower(name) + "%"
		err = db.Find(&tags, "LOWER(name) LIKE ?", nameQuery).Error
	} else {
		err = db.Find(&tags).Error
	}

	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) GetTagById(id string) (*models.Tag, error) {
	var tag models.Tag
	err := database.DB.First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) Update(tag *models.Tag) error {
	err := database.DB.Save(tag).Error
	return err
}

func (r *TagRepository) Delete(tag *models.Tag, isHardDelete bool) error {
	var err error
	if isHardDelete {
		// Delete tag permanently from db
		err = database.DB.Unscoped().Delete(tag).Error
	} else {
		// Soft delete tag
		err = database.DB.Delete(tag).Error
	}
	return err
}
