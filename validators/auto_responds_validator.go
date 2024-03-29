package validators

import (
	"dpacks-go-services-template/models"
	"errors"
	_ "unicode/utf8"
)

func ValidateRespond(autorespond models.AutoRespond, create bool) error {
	if autorespond.Message == "" {
		return errors.New("message cannot be empty")
	}

	if create {
		if autorespond.Message == "" {
			return errors.New("Message cannot be empty")
		}
	}

	return nil
}
