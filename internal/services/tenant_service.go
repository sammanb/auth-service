package services

import (
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
	"github.com/samvibes/vexop/auth-service/internal/utils"
)

type TenantService interface {
	CreateTenant(requester *models.User, name string) (*models.Tenant, error)
	GetTenants(requestor *models.User, page, limit int) ([]*models.Tenant, error)
	GetTenantById(requestor *models.User, id string) (*models.Tenant, error)
	DeleteTenantById(requestor *models.User, id string) (bool, error)
}

type TenantServiceImpl struct {
	repo repository.TenantRepository
}

func NewTenantSvc(repo repository.TenantRepository) TenantService {
	return &TenantServiceImpl{repo}
}

func (s *TenantServiceImpl) CreateTenant(requester *models.User, email string) (*models.Tenant, error) {
	// if requester.Role.Name != utils.RoleSuperAdmin {
	// 	return nil, ErrUnauthorized
	// }

	tenant := &models.Tenant{
		Name:  email,
		Email: email,
	}

	if err := s.repo.CreateTenant(tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *TenantServiceImpl) GetTenants(requestor *models.User, page, limit int) ([]*models.Tenant, error) {
	if requestor.Role.Name != utils.RoleSuperAdmin {
		return nil, ErrUnauthorized
	}

	return s.repo.GetTenants(page, limit)
}

func (s *TenantServiceImpl) GetTenantById(requestor *models.User, id string) (*models.Tenant, error) {
	if requestor.Role.Name != utils.RoleSuperAdmin {
		return nil, ErrUnauthorized
	}

	return s.repo.GetTenantById(id)
}

func (s *TenantServiceImpl) DeleteTenantById(requestor *models.User, id string) (bool, error) {
	if requestor.Role.Name != utils.RoleSuperAdmin {
		return false, ErrUnauthorized
	}

	return s.repo.DeleteTenantById(id)
}
