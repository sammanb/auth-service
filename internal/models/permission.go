package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Action   string    `gorm:"not null" json:"action"`   // CRUD
	Resource string    `gorm:"not null" json:"resource"` // e.g. user, workspace, file
	Code     string    `gorm:"uniqueIndex;" json:"code"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
