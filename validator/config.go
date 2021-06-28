package validator

import (
	cfg "auth2/config"
	"fmt"
)

func ValidateConfig(config *cfg.Config) error {
	if len(config.ListenAddress) < 1 {
		return fmt.Errorf("server address not specified")
	}

	if len(config.DBConnection) < 1 {
		return fmt.Errorf("database connection not specified")
	}

	if len(config.UsersTable) < 1 {
		return fmt.Errorf("users table name not specified")
	}

	return nil
}
