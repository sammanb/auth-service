package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID       uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_tenant_code" json:"tenant_id"`
	Action   string     `gorm:"not null" json:"action"`   // CRUD
	Resource string     `gorm:"not null" json:"resource"` // e.g. user, workspace, file
	Code     string     `gorm:"not null;uniqueIndex:idx_tenant_code" json:"code"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
