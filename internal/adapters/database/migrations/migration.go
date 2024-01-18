package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/adapters/database"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/models"
)

var migrations = []*gormigrate.Migration{}

func Migrate() {
	m := gormigrate.New(database.DB, gormigrate.DefaultOptions, migrations)

	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			&models.Tag{},
		)
		if err != nil {
			return err
		}

		return nil
	})

	// Run the migrations
	if err := m.Migrate(); err != nil {
		panic(err)
	}
}
