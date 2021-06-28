package models

const (
	RoleUser       = "ROLE_USER"
	RoleSuperAdmin = "ROLE_SUPER_ADMIN"
)

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
