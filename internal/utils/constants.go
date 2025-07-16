package utils

import "net/http"

const UserContextKey = "currentUser"

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

var MethodToAction = map[string]string{
	http.MethodGet:    string(ActionRead),
	http.MethodPost:   string(ActionCreate),
	http.MethodPut:    string(ActionUpdate),
	http.MethodPatch:  string(ActionUpdate),
	http.MethodDelete: string(ActionDelete),
}
