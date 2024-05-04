package models

import "time"

type TemplateModel struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	MainFile        string    `json:"mainfile"`
	ThmbnlFile      string    `json:"thmbnlfile"`
	UserID          int       `json:"userid"`
	DevpDescription string    `json:"dmessage"`
	Price           float64   `json:"price"`
	Sdate           time.Time `json:"submitteddate"`
	Status          int       `json:"status"`
	Rating          float64   `json:"rate"`
}
