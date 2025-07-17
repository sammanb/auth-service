package app

import (
	"github.com/samvibes/vexop/auth-service/config"
	"github.com/samvibes/vexop/auth-service/internal/handlers"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"github.com/samvibes/vexop/auth-service/seed"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB            *gorm.DB
	AuthHandler   *handlers.AuthHandler
	TenantHandler *handlers.TenantHandler
	InviteHandler *handlers.InviteHandler
}

func InitApp() *AppContainer {
	db := config.InitDB()
	db.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.Role{},
		&models.Permission{},
		&models.Invitation{},
	)

	seed.SeedSuperAdmin(db)
	seed.SeedRoles(db)

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(*userService)

	tenantRepo := repository.NewTenantRepo(db)
	tenantService := services.NewTenantSvc(tenantRepo)
	tenantHandler := handlers.NewTenantHandler(tenantService)

	inviteRepo := repository.NewInviteRepository(db)
	inviteService := services.NewInviteService(inviteRepo, userRepo)
	inviteHandler := handlers.NewInviteHandler(*inviteService)

	return &AppContainer{
		DB:            db,
		AuthHandler:   authHandler,
		TenantHandler: tenantHandler,
		InviteHandler: inviteHandler,
	}
}
