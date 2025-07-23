package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID    *uuid.UUID    `gorm:"type:uuid;uniqueIndex:idx_tenant_name" json:"tenant_id"`
	Name        string        `gorm:"uniqueIndex:idx_tenant_name;not null" json:"name"`
	Permissions []*Permission `gorm:"many2many:role_permissions" json:"permissions"`
	IsDefault   bool          `gorm:"default:false" json:"is_default"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
