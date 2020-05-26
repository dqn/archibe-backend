package models

import "time"

type Channel struct {
	ID        int64      `db:"id" json:"id"`
	ChannelID string     `db:"channel_id" json:"channel_id"`
	Name      string     `db:"name" json:"name"`
	ImageURL  string     `db:"image_url" json:"image_url"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
