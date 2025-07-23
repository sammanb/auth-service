package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	roleRepo   repository.RoleRepository
}

func NewInviteService(inviteRepo repository.InviteRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository) *InviteService {
	return &InviteService{inviteRepo: inviteRepo, userRepo: userRepo, roleRepo: roleRepo}
}

func (i *InviteService) CreateInvite(requestor *models.User, email, role string, tenant_id uuid.UUID) (string, error) {

	existing_invite, _ := i.inviteRepo.GetInviteByEmailTenant(email, tenant_id.String())

	if existing_invite != nil {
		err := utils.NewAppError(http.StatusConflict, "invite already exists")
		return "", err
	}

	token, hashedToken, err := utils.GenerateRandomToken()
	if err != nil {
		return "", err
	}

	_, err = i.roleRepo.GetRoleByName(tenant_id.String(), role)
	if err != nil {
		return "", utils.NewAppError(http.StatusBadRequest, "invalid role name")
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

func (i *InviteService) GetInvites(requestor *models.User, page, limit int) ([]*dto.InviteResponse, error) {
	return i.inviteRepo.GetInvites(requestor.TenantID.String(), page, limit)
}

func (i *InviteService) RemoveInvite(invite_id string) error {
	if _, err := uuid.Parse(invite_id); err != nil {
		appError := utils.NewAppError(http.StatusBadRequest, "invalid invite id")
		return appError
	}

	err := i.inviteRepo.RemoveInvite(invite_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		appError := utils.NewAppError(http.StatusNotFound, "invite not found")
		return appError
	}
	return err
}

func (i *InviteService) AcceptInvite(acceptInviteReq dto.AcceptInviteRequest, db *gorm.DB) error {
	email := acceptInviteReq.Email
	password := acceptInviteReq.Password
	token := acceptInviteReq.Token
	tenant_id := acceptInviteReq.TenantID

	var appError *utils.AppError

	invite, err := i.inviteRepo.GetInviteByEmailTenant(email, tenant_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appError = utils.NewAppError(http.StatusNotFound, "no invite found for email")
			return appError
		}
		return err
	}

	if invite == nil {
		appError = utils.NewAppError(http.StatusNotFound, "invite expired or deleted")
		return appError
	}

	if invite.Accepted {
		appError = utils.NewAppError(http.StatusBadRequest, "invite already accepted")
		return appError
	}

	// check if token matches invite's hashedtoken
	hashedToken := invite.TokenHash
	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	if err != nil {
		appError = utils.NewAppError(http.StatusBadRequest, "incorrect token")
		return appError
	}

	user, _ := i.userRepo.FindUserByEmailAndTenant(email, tenant_id)
	if user != nil {
		appError = utils.NewAppError(http.StatusBadRequest, "user already exists")
		return appError
	}

	tenantID, err := uuid.Parse(tenant_id)
	if err != nil {
		return errors.New("failed to read tenant id")
	}

	// get role
	var role models.Role
	if err := db.Where("name = ?", invite.Role).First(&role).Error; err != nil {
		appError = utils.NewAppError(http.StatusBadRequest, "invalid role")
		return appError
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

	err = db.Transaction(func(tx *gorm.DB) error {
		err := i.userRepo.CreateUserTx(tx, user)
		if err != nil {
			return errors.New("failed to create user: " + err.Error())
		}

		err = i.inviteRepo.AcceptInviteTx(tx, invite.ID)
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
	// TODO: send mail
	return nil
}
