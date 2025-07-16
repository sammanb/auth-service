package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"not null" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
