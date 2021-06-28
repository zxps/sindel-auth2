package validator

import (
	"auth2/services/user_service"
	"fmt"
)

func ValidateUserRoleIds(roleIds []string, userService *user_service.Service) error {
	for _, roleId := range roleIds {
		role := userService.GetRole(roleId)
		if role == nil {
			return fmt.Errorf("role %s not exists", roleId)
		}
	}

	return nil
}
