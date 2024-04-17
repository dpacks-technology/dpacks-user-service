package models

// Site struct
type Site struct {
	ID          string `json:"id"`
	SeqID       int    `json:"seq_id"`
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Status      int    `json:"status"`
	LastUpdated string `json:"last_updated"`
}
