package seed

import (
	"fmt"
	"log"

	"github.com/samvibes/vexop/auth-service/internal/models"
	"gorm.io/gorm"
)

// A few example Roles
// 	Admin
// 	Member
// 	Guest
// 	Editor
// 	Viewer
// 	Contributor

type Action string

var (
	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Role string

var (
	Admin  Role = "admin"
	Member Role = "member"
	Guest  Role = "guest"
)

type Resource string

var (
	ResourceUser      Resource = "user"
	ResourceFile      Resource = "file"
	ResourceWorkspace Resource = "workspace"
)

var defaultActions = map[string][]string{
	string(Admin):  {string(ActionRead), string(ActionCreate), string(ActionUpdate), string(ActionDelete)},
	string(Member): {string(ActionRead), string(ActionCreate), string(ActionUpdate)},
	string(Guest):  {string(ActionRead)},
}

var defaultResources = map[string][]string{
	string(Admin):  {string(ResourceFile), string(ResourceWorkspace), string(ResourceUser)},
	string(Member): {string(ResourceFile), string(ResourceWorkspace)},
	string(Guest):  {string(ResourceFile)},
}

func SeedRoles(db *gorm.DB) error {
	adminActions := defaultActions[string(Admin)]
	adminResources := defaultResources[string(Admin)]
	isDefault := true
	err := SetResourcePermissions(db, string(Admin), adminActions, adminResources, isDefault)
	if err != nil {
		log.Println("failed to create admin roles")
	}

	isDefault = false
	memberActions := defaultActions[string(Member)]
	memberResources := defaultActions[string(Member)]

	err = SetResourcePermissions(db, string(Member), memberActions, memberResources, isDefault)
	if err != nil {
		log.Println("failed to create member roles")
	}

	guestActions := defaultActions[string(Guest)]
	guestResources := defaultResources[string(Guest)]
	err = SetResourcePermissions(db, string(Guest), guestActions, guestResources, isDefault)
	if err != nil {
		log.Println("failed to create guest roles")
	}

	return nil
}

func SetResourcePermissions(db *gorm.DB, role string, actions []string, resources []string, isDefault bool) error {
	var permissions []*models.Permission
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, action := range actions {
		for _, resource := range resources {
			code := fmt.Sprintf("%s:%s", resource, action)
			var perm models.Permission
			if err := db.Where("code = ?", code).First(&perm).Error; err == gorm.ErrRecordNotFound {
				permission := &models.Permission{
					Action:   action,
					Resource: resource,
					Code:     code,
				}
				err := tx.Create(permission).Error
				if err != nil {
					log.Printf("failed to create permission %s\n", code)
					return err
				}
				permissions = append(permissions, permission)
			}
		}
	}
	newRole := &models.Role{
		Name:        role,
		Permissions: permissions,
		IsDefault:   isDefault,
	}
	err := tx.Create(newRole).Error
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
