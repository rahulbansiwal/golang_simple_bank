package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-z0-9\s]+$`).MatchString
)

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive int")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}

func ValidateString(s string, minLength int, maxLength int) error {
	n := len(s)
	if n < minLength || n > maxLength {
		return fmt.Errorf("length is long")
	}
	return nil
}

func ValidateUsername(username string) error {
	err := ValidateString(username, 3, 10)
	if err != nil {
		return err
	}
	if !isValidUsername(username) {
		return fmt.Errorf("not a valid username")
	}
	return nil
}

func ValidatePassword(password string) error {
	err := ValidateString(password, 6, 100)
	if err != nil {
		return err
	}
	return nil
}

func ValidateEmail(email string) error {
	err := ValidateString(email, 6, 100)
	if err != nil {
		return err
	}
	if _, err = mail.ParseAddress(email); err != nil {
		return fmt.Errorf("not a valid email address")
	}
	return nil
}

func ValidateFullName(username string) error {
	err := ValidateString(username, 3, 10)
	if err != nil {
		return err
	}
	if !isValidFullName(username) {
		return fmt.Errorf("not a valid full name")
	}
	return nil
}
