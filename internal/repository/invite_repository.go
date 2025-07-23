package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"gorm.io/gorm"
)

type InviteRepository interface {
	CreateInvite(*models.Invitation) error
	GetInvites(string, int, int) ([]*dto.InviteResponse, error)
	GetInviteById(string) (*models.Invitation, error)
	GetInviteByEmailTenant(string, string) (*models.Invitation, error)
	RemoveInvite(string) error
	AcceptInviteTx(tx *gorm.DB, inviteID uuid.UUID) error
}

type InviteRepo struct {
	db *gorm.DB
}

func NewInviteRepository(db *gorm.DB) InviteRepository {
	return &InviteRepo{db: db}
}

func (i *InviteRepo) CreateInvite(invitation *models.Invitation) error {
	if err := i.db.Create(invitation).Error; err != nil {
		return err
	}
	return nil
}

func (i *InviteRepo) GetInvites(tenant_id string, page, limit int) ([]*dto.InviteResponse, error) {
	var invitations []*models.Invitation
	offset := (page - 1) * limit
	if err := i.db.Preload("Creator").Where("tenant_id = ? AND accepted = ? AND expires_at > ?", tenant_id, false, time.Now()).Find(&invitations).Offset(offset).Limit(limit).Error; err != nil {
		return nil, err
	}

	var result []*dto.InviteResponse

	for _, invite := range invitations {
		creator := &dto.CreatorInfo{
			Email: invite.Creator.Email,
		}
		_invite := &dto.InviteResponse{
			ID:        invite.ID,
			Email:     invite.Email,
			TenantID:  invite.TenantID,
			Role:      invite.Role,
			Creator:   *creator,
			ExpiresAt: invite.ExpiresAt,
			Accepted:  invite.Accepted,
			CreatedAt: invite.CreatedAt,
			UpdatedAt: invite.UpdatedAt,
		}
		result = append(result, _invite)
	}

	return result, nil
}

func (i *InviteRepo) GetInviteById(invite_id string) (*models.Invitation, error) {
	var invitation models.Invitation
	if err := i.db.Where("id = ?", invite_id).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (i *InviteRepo) GetInviteByEmailTenant(email, tenant_id string) (*models.Invitation, error) {
	var invitation models.Invitation
	if err := i.db.Where("email = ? AND tenant_id = ? AND expires_at > ?", email, tenant_id, time.Now()).First(&invitation).Error; err != nil {
		return nil, err
	}

	return &invitation, nil
}

func (i *InviteRepo) RemoveInvite(inviteID string) error {
	invite := models.Invitation{}

	if err := i.db.Where("id = ?", inviteID).First(&invite).Error; err != nil {
		return err
	}

	if err := i.db.Delete(&invite).Error; err != nil {
		return err
	}

	return nil
}

func (i *InviteRepo) AcceptInviteTx(tx *gorm.DB, inviteID uuid.UUID) error {
	return tx.Model(&models.Invitation{}).Where("id = ?", inviteID).Update("accepted", true).Error
}
