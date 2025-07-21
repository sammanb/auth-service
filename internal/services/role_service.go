package services

import (
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/repository"
)

type RoleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (r *RoleService) GetRoles(tenant_id string, page, limit int) ([]*dto.RoleResponse, error) {
	return r.roleRepo.GetRoles(tenant_id, page, limit)
}

func (r *RoleService) AddRole(tenant_id, name string) error {
	err := r.roleRepo.AddRole(tenant_id, name)
	if err != nil {
		return err
	}

	return nil
}
