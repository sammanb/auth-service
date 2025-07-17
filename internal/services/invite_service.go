package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/config"
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type InviteServiceInterface interface {
	CreateInvite(string, string, string) error
	RemoveInvite(string) error
	AcceptInvite(dto.AcceptInviteRequest) error
	ResendInvite(string) error
}

type InviteService struct {
	inviteRepo repository.InviteRepository
	userRepo   repository.UserRepository
}

func NewInviteService(inviteRepo repository.InviteRepository, userRepo repository.UserRepository) *InviteService {
	return &InviteService{inviteRepo: inviteRepo, userRepo: userRepo}
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
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedBy: requestor.ID,
		Creator:   *requestor,
	}
	err = i.inviteRepo.CreateInvite(invite)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (i *InviteService) GetInvites(requestor *models.User, page, limit int) ([]*models.Invitation, error) {
	return i.inviteRepo.GetInvites(requestor.TenantID.String(), page, limit)
}

func (i *InviteService) RemoveInvite(invite_id string) error {
	return i.inviteRepo.RemoveInvite(invite_id)
}

func (i *InviteService) AcceptInvite(acceptInviteReq dto.AcceptInviteRequest) error {
	email := acceptInviteReq.Email
	password := acceptInviteReq.Password
	token := acceptInviteReq.Token
	tenant_id := acceptInviteReq.TenantID

	invite, err := i.inviteRepo.GetInviteByEmailTenant(email, tenant_id)
	if err != nil {
		return err
	}

	if invite == nil {
		return errors.New("invite expired or deleted")
	}

	if invite.Accepted {
		return errors.New("invite already accepted")
	}

	// check if token matches invite's hashedtoken
	hashedToken := invite.TokenHash
	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	if err != nil {
		return errors.New("incorrect token")
	}

	user, _ := i.userRepo.FindUserByEmailAndTenant(email, tenant_id)
	if user != nil {
		return errors.New("user already exists")
	}

	tenantID, err := uuid.Parse(tenant_id)
	if err != nil {
		return errors.New("failed to read tenant id")
	}

	// get role
	var role models.Role
	if err := config.DB.Where("name = ?", invite.Role).First(&role).Error; err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	user = &models.User{
		TenantID:     &tenantID,
		Email:        email,
		RoleID:       role.ID.String(),
		Role:         role,
		PasswordHash: string(hashedPassword),
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		err := i.userRepo.CreateUserTx(tx, user)
		if err != nil {
			return errors.New("failed to create user: " + err.Error())
		}

		err = i.inviteRepo.AcceptInvite(tx, invite.ID)
		if err != nil {
			return errors.New("failed to accept invite: " + err.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (i *InviteService) ResendInvite(invite_id string) error {
	_, err := i.inviteRepo.GetInviteById(invite_id)
	if err != nil {
		return err
	}
	// send mail
	return nil
}
