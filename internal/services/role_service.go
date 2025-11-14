package services

import (
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/repository"
)

type RoleService interface {
	GetRoles(tenant_id string, page, limit int) ([]*dto.RoleResponse, error)
	AddRole(tenant_id, name string) error
	DeleteRole(id string) error
}

type RoleServiceImpl struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &RoleServiceImpl{roleRepo: roleRepo}
}

func (r *RoleServiceImpl) GetRoles(tenant_id string, page, limit int) ([]*dto.RoleResponse, error) {
	return r.roleRepo.GetRoles(tenant_id, page, limit)
}

func (r *RoleServiceImpl) AddRole(tenant_id, name string) error {
	err := r.roleRepo.AddRole(tenant_id, name)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoleServiceImpl) DeleteRole(id string) error {
	return r.roleRepo.DeleteRole(id)
}
