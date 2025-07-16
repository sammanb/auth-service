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
	IsDefault   bool          `gorm:"type:boolean;default:false" json:"is_default"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
