package main

import "strings"

type PasswordValidationError struct {
	message string
}

func (p *PasswordValidationError) Error() string {
	return p.message
}

func validatePassword(pass string) error {
	if len(pass) < 8 {
		return &PasswordValidationError{"The password must be at least 8 characters long"}
	}
	if strings.ToLower(pass) == pass {
		return &PasswordValidationError{"The password must contain uppercase and lowercase letters"}
	}
	return nil
}
