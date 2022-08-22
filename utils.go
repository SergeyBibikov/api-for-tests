package main

import (
	"regexp"
	"strings"
)

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
	re := regexp.MustCompile(`^([a-zA-Z0-9\.\-\+_]+)@([a-zA-Z0-9\.\-_]+)\.([a-zA-Z]{2,5})$`)

	const message = "The email has an invalid format"
	if isFound := re.Find([]byte(email)); isFound == nil {
		return &EmailValidationError{message}
	}
	return nil
}
