package services

import (
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
)

type InviteServiceInterface interface {
	CreateInvite(string, string, string) error
	RemoveInvite(string) error
	AcceptInvite(string, string) error
}

type InviteService struct {
	inviteRepo repository.InviteRepository
}

func NewInviteService(inviteRepo repository.InviteRepository) *InviteService {
	return &InviteService{inviteRepo: inviteRepo}
}

func (i *InviteService) CreateInvite(requestor *models.User, email, role string, tenant_id uuid.UUID) error {
	return i.inviteRepo.CreateInvite(email, role, tenant_id)
}
