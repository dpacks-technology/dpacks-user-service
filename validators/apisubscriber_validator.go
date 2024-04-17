package validators

import (
	"dpacks-go-services-template/models"
	"errors"
)

func ValidateUserId(subscriber models.ApiSubscriber, create bool) error {

	if subscriber.UserID == "" {
		return errors.New("UserId cannot be empty")
	}

	return nil
}
