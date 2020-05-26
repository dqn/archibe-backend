package models

import "time"

type Video struct {
	ID        int64     `db:"id"`
	VideoID   string    `db:"video_id"`
	ChannelID string    `db:"channel_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
