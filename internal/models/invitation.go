package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invitation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"not null" json:"email"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	Role      string    `gorm:"not null" json:"role"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null"`
	Accepted  bool      `gorm:"not null;default:false" json:"accepted"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	Creator   User      `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
