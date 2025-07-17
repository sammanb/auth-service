package dto

import "github.com/google/uuid"

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	TenantID string `json:"tenant_id" binding:"required,uuid"`
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
