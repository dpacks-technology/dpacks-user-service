package models

// KeyPairs struct is a row record of the keyPairs table in the postgres database
type KeyPairs struct {
	ID         int    `json:"id"`
	UserID     string `json:"user_id"`
	ClientID   string `json:"client_id"`
	Key        string `json:"key"`
	Created_on string `json:"created_on"`
}
