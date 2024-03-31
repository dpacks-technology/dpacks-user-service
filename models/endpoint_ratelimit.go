package models

type EndpointRateLimit struct {
	Id        int    `json:"id"`
	Path      string `json:"path"`
	Limit     int    `json:"limit"`
	CreatedOn string `json:"created_on"`
}
