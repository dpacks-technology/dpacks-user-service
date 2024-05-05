package models

import "time"

type transaction struct {
	TransactionID    int       `json:"transaction_id"`
	UserID           int       `json:"user_id"`
	Plan_id          int       `json:"plan_idn"`
	Amount           float64   `json:"amount"`
	Transaction_Date time.Time `json:"transaction_date  "`
	Status           int       `json:"status"`
}
