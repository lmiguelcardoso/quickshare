package model

import "time"

type ShortLink struct {
	ID string `json:"id"`
	Slug string `json:"slug"`
	OriginalLink string `json:"original_link"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}