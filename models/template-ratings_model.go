package models

import "time"

type TemplateRatingsModel struct {
	TRID     int       `json:"trid"`
	TID      int       `json:"id"`
	UserID   int       `json:"userid"`
	Rating   int       `json:"rating"`
	RateDate time.Time `json:"ratedate"`
}
