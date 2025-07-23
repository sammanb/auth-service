package services

import (
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
	"github.com/samvibes/vexop/auth-service/internal/utils"
)

type TenantSvcInterface interface {
	CreateTenant(requester *models.User, name string) (*models.Tenant, error)
	GetTenants(requestor *models.User, page, limit int) ([]*models.Tenant, error)
	GetTenantById(requestor *models.User, id string) (*models.Tenant, error)
	DeleteTenantById(requestor *models.User, id string) (bool, error)
}

type TenantSvc struct {
	repo repository.TenantRepository
}

func NewTenantSvc(repo repository.TenantRepository) *TenantSvc {
	return &TenantSvc{repo}
}

func (s *TenantSvc) CreateTenant(requester *models.User, name string) (*models.Tenant, error) {
	if requester.Role.Name != utils.RoleSuperAdmin {
		return nil, ErrUnauthorized
	}

	tenant := &models.Tenant{
		Name: name,
	}

	if err := s.repo.CreateTenant(tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *TenantSvc) GetTenants(requestor *models.User, page, limit int) ([]*models.Tenant, error) {
	if requestor.Role.Name != utils.RoleSuperAdmin {
		return nil, ErrUnauthorized
	}

	return s.repo.GetTenants(page, limit)
}

func (s *TenantSvc) GetTenantById(requestor *models.User, id string) (*models.Tenant, error) {
	if requestor.Role.Name != utils.RoleSuperAdmin {
		return nil, ErrUnauthorized
	}

	return s.repo.GetTenantById(id)
}

func (s *TenantSvc) DeleteTenantById(requestor *models.User, id string) (bool, error) {
	if requestor.Role.Name != utils.RoleSuperAdmin {
		return false, ErrUnauthorized
	}

	return s.repo.DeleteTenantById(id)
}
