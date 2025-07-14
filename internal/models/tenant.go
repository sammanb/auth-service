package models

import "gorm.io/gorm"

type Tenant struct {
	ID   string `gorm:"type:uuid;primary_key;" json:"id"`
	Name string `gorm:"not null" json:"name"`
	gorm.Model
}
