package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var RoleSuperAdmin = "superadmin"

type Role struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string        `gorm:"uniqueIndex;not null" json:"name"`
	Permissions []*Permission `gorm:"many2many:role_permissions" json:"permissions"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
