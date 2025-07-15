package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"not null" json:"name"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
