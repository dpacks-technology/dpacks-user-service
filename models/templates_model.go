package models

type Template struct {
	TemplId              int    `json:"templId"`
	TemplName            string `json:"templName"`
	TemplDescription     string `json:"templDescription"`
	TemplCategory        string `json:"templCategory"`
	MainTemplFile        string `json:"mainTemplFile"`
	ThmbnlTemplFile      string `json:"thmbnlTemplFile"`
	TemplDevpName        string `json:"templDevpName"`
	UserID               int    `json:"userID"`
	TemplDevpDescription string `json:"templDevpDescription"`
}
