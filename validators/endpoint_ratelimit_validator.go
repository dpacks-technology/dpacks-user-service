package validators

import (
	"dpacks-go-services-template/models"
	"errors"
	"strings"
)

func ValidatePath(ratelimit models.EndpointRateLimit, create bool) error {

	if create {
		if ratelimit.Path == "" {
			return errors.New("path cannot be empty")
		}
		if ratelimit.Limit == 0 {
			return errors.New("Limit cannot be empty")
		}
		if !strings.HasPrefix(ratelimit.Path, "/") {
			return errors.New("Path must start with a forward slash (/)")
		}
	}

	return nil
}
