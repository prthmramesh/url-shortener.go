package main

import "time"

type ShortURL struct {
	ShortCode   string    `json:"shortCode"`
	OriginalURL string    `json:"originalURL"`
	CreatedTime time.Time `json:"timestamp"`
}

type ClickEvent struct {
	ShortCode string
	Timestamp time.Time
}
