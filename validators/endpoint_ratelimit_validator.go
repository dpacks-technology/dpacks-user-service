package validators

import (
	"dpacks-go-services-template/models"
	"errors"
	"strings"
)

func ValidatePath(endpoint models.Endpoint, create bool) error {

	if create {
		if endpoint.Path == "" {
			return errors.New("path cannot be empty")
		}
		if endpoint.Limit == 0 {
			return errors.New("Limit cannot be empty")
		}
		if !strings.HasPrefix(endpoint.Path, "/") {
			return errors.New("Path must start with a forward slash (/)")
		}
	}

	return nil
}
