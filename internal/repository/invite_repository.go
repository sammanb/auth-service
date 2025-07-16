package repository

import (
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"gorm.io/gorm"
)

type InviteRepository interface {
	CreateInvite(string, string, uuid.UUID) error
	GetInvites(string, int, int) (*[]models.Invitation, error)
	RemoveInvite(string) error
	AcceptInvite(string) error
}

type InviteRepo struct {
	db *gorm.DB
}

func NewInviteRepository(db *gorm.DB) InviteRepository {
	return &InviteRepo{db: db}
}

func (i *InviteRepo) CreateInvite(email, role string, tenantId uuid.UUID) error {
	invite := models.Invitation{
		Email:    email,
		TenantID: tenantId,
		Role:     role,
	}
	if err := i.db.Create(invite).Error; err != nil {
		return err
	}
	return nil
}

func (i *InviteRepo) GetInvites(tenant_id string, page, limit int) (*[]models.Invitation, error) {
	var invitations []models.Invitation
	if err := i.db.Where("tenant_id = ?", tenant_id).Find(&invitations).Offset(page).Limit(limit).Error; err != nil {
		return nil, err
	}
	return &invitations, nil
}

func (i *InviteRepo) RemoveInvite(inviteID string) error {
	invite := models.Invitation{}
	// check if invite exists
	if err := i.db.Where("id = ?", inviteID).First(&invite).Error; err != nil {
		return err
	}
	if err := i.db.Delete(invite).Error; err != nil {
		return err
	}

	return nil
}

func (i *InviteRepo) AcceptInvite(inviteID string) error {
	return nil
}
