package validators

import (
	"dpacks-go-services-template/models"
	"errors"
)

func ValidateUserId(kepair models.KeyPairs, create bool) error {

	if kepair.UserID == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}
