package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
	"github.com/samvibes/vexop/auth-service/internal/utils"
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

func (i *InviteService) CreateInvite(requestor *models.User, email, role string, tenant_id uuid.UUID) (string, error) {
	token, hashedToken, err := utils.GenerateInviteToken()
	if err != nil {
		return "", err
	}
	invite := &models.Invitation{
		Email:     email,
		TenantID:  tenant_id,
		Role:      role,
		TokenHash: hashedToken,
		ExpiresAt: time.Now(),
		CreatedBy: requestor.ID,
		Creator:   *requestor,
	}
	err = i.inviteRepo.CreateInvite(invite)
	if err != nil {
		return "", err
	}

	return token, nil
}
