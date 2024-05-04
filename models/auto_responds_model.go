package models

type AutoRespond struct {
	ID          int    `json:"id"`
	Message     string `json:"message"`
	Trigger     string `json:"trigger"`
	LastUpdated string `json:"last_updated"`
	Status      int    `json:"status"`
	WebID       string `json:"webid"`
}
