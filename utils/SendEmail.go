package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Email struct {
	ApiKey  string `json:"api_key"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Size    string `json:"size"`
}

func SendEmail(to, subject, message, size string) error {
	email := Email{
		ApiKey:  os.Getenv("EMAIL_API_KEY"),
		To:      to,
		Subject: subject,
		Message: message,
		Size:    size,
	}

	emailJson, err := json.Marshal(email)
	if err != nil {
		return err
	}

	// Replace with the URL of your email microservice
	url := os.Getenv("EMAIL_API_HOST") + "/api/email/send"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(emailJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status of the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email: status code %d", resp.StatusCode)
	}

	return nil
}
