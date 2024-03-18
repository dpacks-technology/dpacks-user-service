package models

type WebpageModel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	WebID string `json:"webid"`
	Path  string `json:"path"`
}
