package models

// UserAlerts struct is a row record of the UserAlerts table in the postgres database
type UserAlerts struct {
	AlertID            int    `json:"alert_id"`
	UserID             int    `json:"user_id"`
	UserEmail          string `json:"user_email"`
	AlertThreshold     int    `json:"alert_threshold"`
	AlertSubject       string `json:"alert_subject"`
	AlertContent       string `json:"alert_content"`
	WhenAlertRequired  string `json:"when_alert_required"`
	ReminderOption     string `json:"reminder_option"`
	CustomReminderDate string `json:"custom_reminder_date"`
}
type VisitorInfo struct {
	ID       int    `json:"id"`
	IpAddres string `json:"ip_adrees"`
	Device   string `json:"device"`
	Country  string `json:"country"`
	Source   string `json:"source"`
}
