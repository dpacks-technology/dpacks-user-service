package models

type UserAlertsShow struct {
	AlertID           int    `json:"id" db:"id"`
	AlertThreshold    int    `json:"alert_threshold" db:"alert_threshold"`
	AlertSubject      string `json:"alert_subject" db:"alert_subject"`
	AlertContent      string `json:"alert_content" db:"alert_content"`
	WhenAlertRequired string `json:"when_alert_required" db:"when_alert_required"`
	RepeatOn          string `json:"repeat_on" db:"repeat_on"`
	Status            int    `json:"status" db:"status"`
	WebsiteeId        string `json:"website_id" db:"website_id"`
}

type UserAlertStatus struct {
	Status int `json:"status" db:"status"`
}

type CreateNewUserAlert struct {
	AlertThreshold    int    `json:"Threshold" db:"alert_threshold"`
	AlertSubject      string `json:"Subject" db:"alert_subject"`
	AlertContent      string `json:"Content" db:"alert_content"`
	WhenAlertRequired string `json:"AlertOn" db:"when_alert_required"`
	RepeatOn          string `json:"Repeat" db:"repeat_on"`
	WebsiteeId        string `json:"webId" db:"website_id"`
}
