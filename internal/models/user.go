package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleAdmin      UserRole = "admin"
	RoleMemeber    UserRole = "member"
	RoleGuest      UserRole = "guest"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	TenantID     *uuid.UUID `gorm:"type:uuid"`
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"`
	Role         UserRole   `gorm:"type:varchar(20);not null"`
	gorm.Model
}
