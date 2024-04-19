package models

type UserAlertsShow struct {
	AlertID           int    `json:"id" db:"id"`
	AlertThreshold    int    `json:"alert_threshold" db:"alert_threshold"`
	AlertSubject      string `json:"alert_subject" db:"alert_subject"`
	AlertContent      string `json:"alert_content" db:"alert_content"`
	WhenAlertRequired string `json:"when_alert_required" db:"when_alert_required"`
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
	WebsiteeId        string `json:"webId" db:"website_id"`
}

type UpdateUserAlert struct {
	AlertThreshold    int    `json:"Threshold" db:"alert_threshold"`
	AlertSubject      string `json:"Subject" db:"alert_subject"`
	AlertContent      string `json:"Content" db:"alert_content"`
	WhenAlertRequired string `json:"AlertOn" db:"when_alert_required"`
}

type SessionRecord struct {
	SessionID    string `json:"session_id" db:"sessionid"`
	SessionStart string `json:"session_start" db:"sessionstart"`
	IpAddress    string `json:"ip_address" db:"ipaddress"`
	CountryCode  string `json:"country_code" db:"countrycode"`
	DeviceId     int    `json:"device_id" db:"deviceid"`
	SourceId     int    `json:"source_id" db:"source_id"`
	LandingPage  string `json:"landing_page" db:"landingpage"`
	WebId        string `json:"web_id" db:"web_id"`
}
