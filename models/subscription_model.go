package models

// Site struct
type SubscriptionModel struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	PlanID    int    `json:"plan_id"`
	PlanName  string `json:"plan_name"`
	Amount    string `json:"amount"`
}
