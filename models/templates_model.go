package models

type TemplateModel struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	MainFile        string  `json:"mainfile"`
	ThmbnlFile      string  `json:"thmbnlfile"`
	DevpName        string  `json:"devpname"`
	UserID          int     `json:"userid"`
	DevpDescription string  `json:"dmessage"`
	Price           float64 `json:"price"`
	Sdate           string  `json:"submitteddate"`
	Status          int     `json:"status"`
}
