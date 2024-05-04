package models

// SubscriptionPlans struct
type SubscriptionPlans struct {
	PlanID       int     `json:"plan_id"`
	PlanName     string  `json:"plan_name"`
	Desc         string  `json:"description"`
	Duration     int     `json:"duration"`
	Features     string  `json:"features"`
	MonthlyPrice float64 `json:"monthly_price"`
	AnnualPrice  float64 `json:"annual_price"`
	Status       int     `json:"status"`
}
