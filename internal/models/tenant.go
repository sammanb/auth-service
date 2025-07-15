package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID   string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"not null" json:"name"`
	gorm.Model

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
