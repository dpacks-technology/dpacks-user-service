package models

type EndpointRateLimit struct {
	Id        int    `json:"id"`
	Path      string `json:"path"`
	Limit     int    `json:"ratelimit"`
	CreatedOn string `json:"created_on"`
	Status    int    `json:"status"`
}
