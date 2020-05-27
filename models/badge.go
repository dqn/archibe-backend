package models

import "time"

type Badge struct {
	OwnerChannelID string    `json:"owner_channel_id"`
	LiverChannelID string    `json:"liver_channel_id"`
	BadgeType      string    `json:"badge_type"`
	ImageURL       string    `json:"image_url"`
	Label          string    `json:"label"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
