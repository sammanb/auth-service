package repository

import (
	"github.com/samvibes/vexop/auth-service/internal/models"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(tenant *models.Tenant) error
	GetTenants(page, limit int) ([]*models.Tenant, error)
	GetTenantById(id string) (*models.Tenant, error)
	DeleteTenantById(id string) (bool, error)
}

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (t *TenantRepo) CreateTenant(tenant *models.Tenant) error {
	return t.db.Create(tenant).Error
}

func (t *TenantRepo) GetTenants(page, limit int) ([]*models.Tenant, error) {
	var tenants []*models.Tenant

	offset := (page - 1) * limit
	tx := t.db.Offset(offset).Limit(limit).Find(&tenants)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tenants, nil
}

func (t *TenantRepo) GetTenantById(id string) (*models.Tenant, error) {
	var tenant models.Tenant
	tx := t.db.Where("id = ?", id).First(&tenant)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &tenant, nil
}

func (t *TenantRepo) DeleteTenantById(id string) (bool, error) {
	tenant, err := t.GetTenantById(id)
	if err != nil {
		return false, err
	}

	tx := t.db.Delete(tenant)
	if tx.Error != nil {
		return false, tx.Error
	}

	return true, nil
}
