package models

// VisitorUser struct is a row record of the visitor_users table in the postgres database
type VisitorUser struct {
	UserID             int    `json:"UserID"`
	Name               string `json:"Name"`
	Email              string `json:"Email"`
	PhoneNumber        string `json:"PhoneNumber"`
	DateOfBirth        string `json:"DateOfBirth"`
	Country            string `json:"Country"`
	FavoriteCategories string `json:"FavoriteCategories"`
	UserDescription    string `json:"UserDescription"`
	SignUpDate         string `json:"SignUpDate"`
	LastLogin          string `json:"LastLogin"`
	ProfilePicture     string `json:"ProfilePicture"`
	Gender             string `json:"Gender"`
	Language           string `json:"Language"`
	Timezone           string `json:"Timezone"`
}
