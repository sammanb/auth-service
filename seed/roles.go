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

var defaultActions = map[string][]string{
	string(utils.Admin):  {string(utils.ActionRead), string(utils.ActionCreate), string(utils.ActionUpdate), string(utils.ActionDelete)},
	string(utils.Member): {string(utils.ActionRead), string(utils.ActionCreate), string(utils.ActionUpdate)},
	string(utils.Guest):  {string(utils.ActionRead)},
}

var defaultResources = map[string][]string{
	string(utils.Admin):  {string(utils.ResourceFile), string(utils.ResourceWorkspace), string(utils.ResourceUser)},
	string(utils.Member): {string(utils.ResourceFile), string(utils.ResourceWorkspace)},
	string(utils.Guest):  {string(utils.ResourceFile)},
}

func SeedRoles(db *gorm.DB) error {
	adminActions := defaultActions[string(utils.Admin)]
	adminResources := defaultResources[string(utils.Admin)]

	err := SetResourcePermissions(db, string(utils.Admin), adminActions, adminResources)
	if err != nil {
		log.Println("failed to create admin roles")
	}

	memberActions := defaultActions[string(utils.Member)]
	memberResources := defaultActions[string(utils.Member)]

	err = SetResourcePermissions(db, string(utils.Member), memberActions, memberResources)
	if err != nil {
		log.Println("failed to create member roles")
	}

	guestActions := defaultActions[string(utils.Guest)]
	guestResources := defaultResources[string(utils.Guest)]
	err = SetResourcePermissions(db, string(utils.Guest), guestActions, guestResources)
	if err != nil {
		log.Println("failed to create guest roles")
	}

	return nil
}

func SetResourcePermissions(db *gorm.DB, role string, actions []string, resources []string) error {
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
