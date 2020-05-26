package models

import "time"

type Channel struct {
	id        int64     `db:"id"`
	ChannelID string    `db:"channel_id"`
	Name      string    `db:"name"`
	ImageURL  string    `db:"image_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
