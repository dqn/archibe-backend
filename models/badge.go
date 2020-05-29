package models

import "time"

type Badge struct {
	ID             int64      `db:"id" json:"id"`
	OwnerChannelID string     `db:"owner_channel_id" json:"owner_channel_id"`
	LiverChannelID string     `db:"liver_channel_id" json:"liver_channel_id"`
	BadgeType      string     `db:"badge_type" json:"badge_type"`
	ImageURL       string     `db:"image_url" json:"image_url"`
	Label          string     `db:"label" json:"label"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
}

type Badges []Badge

func (m *Badges) Scan(val interface{}) error {
	return scanJSON(m, val)
}
