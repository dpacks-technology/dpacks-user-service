package validators

import (
	"dpacks-go-services-template/models"
	"errors"
	"unicode/utf8"
)

func ValidateName(webpage models.WebpageModel, create bool) error {

	if webpage.Name == "" {
		return errors.New("name cannot be empty")
	}
	if utf8.RuneCountInString(webpage.Name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}
	if utf8.RuneCountInString(webpage.Name) > 30 {
		return errors.New("name cannot exceed 30 characters")
	}

	if create {
		if webpage.Path == "" {
			return errors.New("path cannot be empty")
		}
		if webpage.WebID == "" {
			return errors.New("webid cannot be empty")
		}
	}

	return nil
}
