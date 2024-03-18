package models

type AutoRespond struct {
	ID          int    `json:"id"`
	Message     string `json:"message"`
	Trigger     string `json:"trigger"`
	IsActive    bool   `json:"is_active"`
	LastUpdated string `json:"last_updated"`
}
