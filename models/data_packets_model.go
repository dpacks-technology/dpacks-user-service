package models

import "time"

type DataPackets struct {
	ID           string    `json:"id"`
	Site         string    `json:"site"`
	Page         string    `json:"page"`
	Element      string    `json:"element"`
	InitDatetime time.Time `json:"init_dateTime"`
	LastUpdated  time.Time `json:"last_updated"`
	Size         int       `json:"size"`
	Status       int       `json:"status"`
	Pinned       int       `json:"pinned"`
}
