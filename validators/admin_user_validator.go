package validators

import (
	"dpacks-go-services-template/models"
	"errors"
	"unicode/utf8"
)

func ValidateAdmin(admin models.AdminUserModel, create bool) error {

	if admin.Name == "" {
		return errors.New("name cannot be empty")
	}
	if utf8.RuneCountInString(admin.Name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}
	if utf8.RuneCountInString(admin.Name) > 30 {
		return errors.New("name cannot exceed 30 characters")
	}

	if create {
		if admin.Phone == "" {
			return errors.New("phone cannot be empty")
		}
		if admin.Email == "" {
			return errors.New("email cannot be empty")
		}
		if admin.Password == "" {
			return errors.New("password cannot be empty")
		}
	}

	return nil
}
