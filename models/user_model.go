package models

type UserModel struct {
	ID               int     `json:"id"`
	Code             *string `json:"code"`
	DOB              string  `json:"dob"`
	FirstName        string  `json:"first_name"`
	ForgotCode       *string `json:"forgot_code"`
	ForgotCodeExpire *string `json:"forgot_code_expire"`
	Gender           string  `json:"gender"`
	InitDate         string  `json:"init_date"`
	LastName         string  `json:"last_name"`
	Password         string  `json:"password"`
	Phone            string  `json:"phone"`
	Status           int     `json:"status"`
	UserKey          string  `json:"user_key"` //uuid
	Email            string  `json:"email"`
	VerificationExp  *string `json:"verification_exp"`

	//CreatedOn string `json:"created_on"`
}
