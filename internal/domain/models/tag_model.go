package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID uuid.UUID `gorm:"type:uuid;column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	/* Fields */
	Name string `gorm:"not null;uniqueIndex:unique_tag_name" json:"name"`
	/* Timestamp */
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName is Database TableName of this model
func (e *Tag) TableName() string {
	return "tags"
}
