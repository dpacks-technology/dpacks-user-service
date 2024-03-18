package models

// SubscriptionPlans struct
type SubscriptionPlans struct {
	PlanID   int     `json:"plan_id"`
	PlanName string  `json:"plan_name"`
	Desc     string  `json:"description"`
	Price    float64 `json:"price"`
	Duration int     `json:"duration"`
	Features string  `json:"features"`
}
