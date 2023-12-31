package validation

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"

	"gopkg.in/validator.v2"
)

func SetupValidators() {
	validators := []struct {
		name string
		fn   validator.ValidationFunc
	}{
		{"password", password},
		{"email", email},
		{"role", role},
		{"sort", sort},
	}

	for _, v := range validators {
		_ = validator.SetValidationFunc(v.name, v.fn)
	}
}

var (
	passwordMinLength       = 8
	emailMaxLength          = 127
	passwordSpecialChars    = "!$#()[]{}?+*~@^&-_"
	passwordRequiredEntries = []struct {
		name  string
		chars string
	}{
		{"lowercase character", "abcdefghijklmnopqrstuvwxyz"},
		{"uppercase character", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"digit", "0123456789"},
		{"special character (" + passwordSpecialChars + ")", passwordSpecialChars},
	}
)

func password(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("password only validates string")
	}

	if len(s) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters long", passwordMinLength)
	}

	for _, required := range passwordRequiredEntries {
		if !strings.ContainsAny(s, required.chars) {
			return fmt.Errorf("password must contain at least one %s", required.name)
		}
	}
	return nil
}

func email(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("email only validates string")
	}

	if len(s) > emailMaxLength {
		return fmt.Errorf("email must be at least %d characters long", emailMaxLength)
	}
	_, err := mail.ParseAddress(s)
	return err
}

func role(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("email only validates string")
	}

	if !(s == users.AdminRole || s == users.EditorRole || s == users.UserRole) {
		return fmt.Errorf("you can't create new roles: %s ", s)
	}

	return nil
}

func sort(v interface{}, _ string) error {
	validate := func(s *string) error {
		if s == nil {
			return nil
		}
		switch *s {
		case "asc", "desc":
			return nil
		}
		return fmt.Errorf("sort must be one of asc or desc")
	}

	switch s := v.(type) {
	case string:
		return validate(&s)
	case *string:
		return validate(s)
	default:
		return fmt.Errorf("sort only validates strings or pointers to strings")

	}
}
