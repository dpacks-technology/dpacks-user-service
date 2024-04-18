package validators

import (
	"dpacks-go-services-template/models"
	"errors"
)

func ValidateMessage(autorespond models.AutoRespond, create bool) error {

	if autorespond.Message == "" {
		return errors.New("name cannot be empty")
	}

	if create {
		if autorespond.Message == "" {
			return errors.New("path cannot be empty")
		}

	}

	return nil
}
