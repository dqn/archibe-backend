package models

import "time"

type Video struct {
	ID        int64      `db:"id" json:"id"`
	VideoID   string     `db:"video_id" json:"video_id"`
	ChannelID string     `db:"channel_id" json:"channel_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

type Videos []Video

func (v *Videos) Scan(val interface{}) error {
	return scanJSON(v, val)
}
