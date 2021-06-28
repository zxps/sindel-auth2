package validator

import (
	"fmt"
)

type PasswordValidateOptions struct {
	MinLength int
	MaxLength int
}

func PasswordValidate(password string, opts *PasswordValidateOptions) error {
	if len(password) < opts.MinLength {
		return fmt.Errorf("too short. min password length - %s", opts.MinLength)
	}

	if len(password) > opts.MaxLength {
		return fmt.Errorf("too long. max password length - %d", opts.MaxLength)
	}

	return nil
}
