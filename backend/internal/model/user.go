package model

type Role string

const (
	RoleAdminGateway   Role = "admin_gateway"
	RoleAppUser        Role = "app_user"
	RoleMonitoringUser Role = "monitoring_user"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         Role
	AppName      string
}
