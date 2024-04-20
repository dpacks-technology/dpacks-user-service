package validators

import (
	"dpacks-go-services-template/models"
	"errors"
)

func ValidateSubscription(subscription models.SubscriptionModel, create bool) error {

	if subscription.ProjectID == "" {
		return errors.New("project_id cannot be empty")
	}
	if subscription.PlanID == 0 {
		return errors.New("plan_id cannot be empty")
	}

	if create {

	}

	return nil
}
