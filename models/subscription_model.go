package models

type SubscriptionModel struct {
	ProjectID string `json:"project_id"`
	PlanID    int    `json:"plan_id"`
}
