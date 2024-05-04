package models

type AuthRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AuthResponseBody struct {
	ID       int        `json:"id"`
	UserKey  string     `json:"userKey"`
	Username string     `json:"username"`
	Status   int        `json:"status"`
	Role     []AuthRole `json:"role"`
}
