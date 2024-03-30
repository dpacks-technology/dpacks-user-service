package models

type TransactionsModel struct {
	TransactionID   int     `json:"id"`
	UserID          int     `json:"user_id"`
	PlanID          int     `json:"plan_id"`
	Amount          float64 `json:"amount"`
	TransactionDate string  `json:"transaction_date"`
	Status          int     `json:"status"`
	PlanName        string  `json:"plan_name"`
	Email           string  `json:"email"`
	Phone           int     `json:"phone"`
}
