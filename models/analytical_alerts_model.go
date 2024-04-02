package models

type UserAlertsModel struct {
	AlertID            int    `json:"id" db:"id"`
	UserID             int    `json:"user_id" db:"user_id"`
	AlertThreshold     int    `json:"alert_threshold" db:"alert_threshold"`
	AlertSubject       string `json:"alert_subject" db:"alert_subject"`
	AlertContent       string `json:"alert_content" db:"alert_content"`
	WhenAlertRequired  string `json:"when_alert_required" db:"when_alert_required"`
	RepeatOn           string `json:"repeat_on" db:"repeat_on"`
	CustomReminderDate string `json:"custom_reminder_date" db:"custom_reminder_date"`
	Status             int    `json:"status" db:"status"`
	WebsiteeId         string `json:"website_id" db:"website_id"`
}
