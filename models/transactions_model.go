package models

import "time"

type TransactionsModel struct {
	TransactionID   int       `json:"id"`
	WebId           string    `json:"web_id"`
	CompanyName     string    `json:"company_name"`
	StreetNo        string    `json:"street_no"`
	City            string    `json:"city"`
	PostalCode      string    `json:"postal_code"`
	Country         string    `json:"country"`
	Email           string    `json:"email"`
	PaymentMethod   string    `json:"payment_method"`
	GivenName       string    `json:"given_name"`
	LastName        string    `json:"last_name"`
	Month           int       `json:"month"`
	Year            int       `json:"year"`
	CVV             int       `json:"cvv"`
	Terms           bool      `json:"terms"`
	TransactionDate time.Time `json:"transaction_date"`
	Status          int       `json:"status"`
	CardNumber      int       `json:"card_number"`
}
