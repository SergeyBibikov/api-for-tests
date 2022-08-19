package main

import "strings"

type PasswordValidationError struct {
	message string
}

func (p *PasswordValidationError) Error() string {
	return p.message
}

type EmailValidationError struct {
	message string
}

func (e *EmailValidationError) Error() string {
	return e.message
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

func validateEmail(email string) error {
	const message = "The email must be in the following format: <>@<>.com"
	if !strings.Contains(email, "@") {
		return &EmailValidationError{message}
	}
	return nil
}
