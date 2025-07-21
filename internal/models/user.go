package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID               *uuid.UUID `gorm:"type:uuid"`
	Email                  string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash           string     `gorm:"not null" json:"-"`
	RoleID                 string     `json:"role_id"`
	Role                   Role       `gorm:"foreignKey:RoleID" json:"role"`
	IsOwner                bool       `gorm:"not null;default:false" json:"is_owner"`
	ResetPasswordTokenHash string     `json:"reset_password_token_hash"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
