package validators

import (
	"dpacks-go-services-template/models"
	"errors"
)

func ValidateNames(transaction models.TransactionsModel, create bool) error {

	if transaction.StreetNo == "" {
		return errors.New("StreetNo cannot be empty")
	}

	if transaction.City == "" {
		return errors.New("City cannot be empty")
	}

	if transaction.PostalCode == "" {
		return errors.New("PostalCode cannot be empty")
	}

	if transaction.Country == "" {
		return errors.New("Country cannot be empty")
	}

	if transaction.Email == "" {
		return errors.New("Email cannot be empty")
	}

	if transaction.Year == 0 {
		return errors.New("Year cannot be empty")
	}

	if transaction.GivenName == "" {
		return errors.New("PaymentMethod cannot be empty")
	}

	if transaction.Month == 0 {
		return errors.New("Month cannot be empty")
	}

	if transaction.CardNumber == 0 {
		return errors.New("CardNumber cannot be empty")
	}

	if transaction.Amount == 0 {
		return errors.New("Amount cannot be empty")
	}

	if transaction.Terms == true {
		return errors.New("Terms cannot be empty")
	}

	return nil
}
