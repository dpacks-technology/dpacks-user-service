package models

type AdminUserModel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    int    `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}