package seeds

import (
	"product-service/internal/adapters/database"
	"product-service/internal/domain/models"
	"product-service/pkg/logger"
)

func SeedData() error {
	tags := []models.Tag{
		{
			Name: "tag1",
		},
		{
			Name: "tag2",
		},
		{
			Name: "tag3",
		},
	}
	if err := createRecords(&tags); err != nil {
		logger.Errorf("Failed to create records: %s", err)
		return err
	}

	return nil
}

func createRecords(data interface{}) error {
	if err := database.DB.Create(data).Error; err != nil {
		logger.Errorf("Failed to create records: %s", err)
		return err
	}

	return nil
}

func IsSeedDataExists() bool {
	var count int64
	database.DB.Model(&models.Tag{}).Count(&count)
	return count > 0
}
