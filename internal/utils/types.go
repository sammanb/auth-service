package utils

import (
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
)

type PermissionMap map[uuid.UUID]*models.Permission
