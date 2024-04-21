package validators

import (
	"dpacks-go-services-template/models"
	"errors"
	"unicode/utf8"
)

func ValidateSite(site models.Site, create bool) error {

	if site.Name == "" {
		return errors.New("name cannot be empty")
	}

	if utf8.RuneCountInString(site.Name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}

	if utf8.RuneCountInString(site.Name) > 30 {
		return errors.New("name cannot exceed 30 characters")
	}

	if site.Description == "" {
		return errors.New("description cannot be empty")
	}

	if utf8.RuneCountInString(site.Description) < 10 {
		return errors.New("description must be at least 10 characters long")
	}

	if utf8.RuneCountInString(site.Description) > 100 {
		return errors.New("description cannot exceed 100 characters")
	}

	if create {
		if site.Name == "" {
			return errors.New("name cannot be empty")
		}

		if utf8.RuneCountInString(site.Name) < 3 {
			return errors.New("name must be at least 3 characters long")
		}

		if utf8.RuneCountInString(site.Name) > 30 {
			return errors.New("name cannot exceed 30 characters")
		}

		if site.Domain == "" {
			return errors.New("domain cannot be empty")
		}

		if site.Description == "" {
			return errors.New("description cannot be empty")
		}

		if utf8.RuneCountInString(site.Description) < 10 {
			return errors.New("description must be at least 3 characters long")
		}

		if utf8.RuneCountInString(site.Description) > 100 {
			return errors.New("description cannot exceed 100 characters")
		}

	}

	return nil
}
