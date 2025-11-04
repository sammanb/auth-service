package dto

import (
	"time"

	"github.com/google/uuid"
)

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	TenantID string `json:"tenant_id" binding:"omitempty,uuid"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateInviteRequest struct {
	Email    string    `json:"email" binding:"required,email"`
	Role     string    `json:"role" binding:"required"`
	TenantID uuid.UUID `json:"tenant_id" binding:"required"`
}

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

type AcceptInviteRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
	TenantID string `json:"tenant_id" binding:"required"`
}

type UpdateUserRoleRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreatorInfo struct {
	Email string `json:"email"`
}

type InviteResponse struct {
	ID        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	TenantID  uuid.UUID   `json:"tenant_id"`
	Role      string      `json:"role"`
	ExpiresAt time.Time   `json:"expires_at"`
	Accepted  bool        `json:"accepted"`
	Creator   CreatorInfo `json:"creator"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type PermissionInfo struct {
	ID       string `json:"id"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Code     string `json:"code"`
}

type RoleResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Permissions []PermissionInfo `json:"permissions"`
	IsDefault   bool             `json:"is_default"`
}
