package seed

import (
	"fmt"
	"log"

	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"gorm.io/gorm"
)

// A few example Roles
// 	Admin
// 	Member
// 	Guest
// 	Editor
// 	Viewer
// 	Contributor

// All resource permissions
var adminRole = map[utils.Resource][]utils.Action{
	utils.ResourceFile:      {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate, utils.ActionDelete},
	utils.ResourceWorkspace: {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate, utils.ActionDelete},
	utils.ResourceUser:      {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate, utils.ActionDelete},
	utils.ResourceInvite:    {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate, utils.ActionDelete},
	utils.ResourceRole:      {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate, utils.ActionDelete},
}

var memberRole = map[utils.Resource][]utils.Action{
	utils.ResourceFile:      {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate},
	utils.ResourceWorkspace: {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate},
	utils.ResourceUser:      {utils.ActionRead},
	utils.ResourceInvite:    {utils.ActionRead, utils.ActionUpdate},
	utils.ResourceRole:      {utils.ActionRead},
}

var guestRole = map[utils.Resource][]utils.Action{
	utils.ResourceFile:      {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate},
	utils.ResourceWorkspace: {utils.ActionRead, utils.ActionCreate, utils.ActionUpdate},
	utils.ResourceUser:      {utils.ActionRead},
	utils.ResourceInvite:    {utils.ActionRead, utils.ActionUpdate},
	utils.ResourceRole:      {},
}

func SeedRoles(db *gorm.DB) error {
	permissions := make([]*models.Permission, 0)
	permissionMap := make(map[string]*models.Permission, 0)
	for resource := range adminRole {
		actions := adminRole[resource]
		_permissions, err := CreatePermissions(db, resource, actions)
		if err != nil {
			log.Printf("failed to create admin role for %s\n", resource)
		}
		permissions = append(permissions, _permissions...)
		for _, perm := range permissions {
			permissionMap[perm.Code] = perm
		}
	}
	CreateRole(db, string(utils.Admin), permissions, true)

	permissions = make([]*models.Permission, 0)
	for resource := range memberRole {
		actions := memberRole[resource]
		for _, action := range actions {
			code := fmt.Sprintf("%s:%s", resource, action)
			perm := permissionMap[code]
			permissions = append(permissions, perm)
		}
	}
	CreateRole(db, string(utils.Member), permissions, false)

	permissions = make([]*models.Permission, 0)
	for resource := range guestRole {
		actions := guestRole[resource]
		for _, action := range actions {
			code := fmt.Sprintf("%s:%s", resource, action)
			perm := permissionMap[code]
			permissions = append(permissions, perm)
		}
	}
	CreateRole(db, string(utils.Guest), permissions, false)

	return nil
}

func CreatePermissions(db *gorm.DB, resource utils.Resource, actions []utils.Action) ([]*models.Permission, error) {
	var permissions []*models.Permission
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, action := range actions {
		code := fmt.Sprintf("%s:%s", resource, action)
		var perm models.Permission
		if err := db.Where("code = ?", code).First(&perm).Error; err == gorm.ErrRecordNotFound {
			permission := &models.Permission{
				Action:   string(action),
				Resource: string(resource),
				Code:     code,
			}
			err := tx.Create(permission).Error
			if err != nil {
				log.Printf("failed to create permission %s\n", code)
				return nil, err
			}
			permissions = append(permissions, permission)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func CreateRole(db *gorm.DB, roleName string, permissions []*models.Permission, isDefault bool) error {
	newRole := &models.Role{
		Name:        roleName,
		Permissions: permissions,
		IsDefault:   isDefault,
	}
	err := db.Create(newRole).Error
	if err != nil {
		return err
	}

	return nil
}
