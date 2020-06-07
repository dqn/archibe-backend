package models

import "time"

type Video struct {
	ID            int64      `db:"id" json:"id"`
	VideoID       string     `db:"video_id" json:"video_id"`
	ChannelID     string     `db:"channel_id" json:"channel_id"`
	Title         string     `db:"title" json:"title"`
	Description   string     `db:"description" json:"description"`
	LengthSeconds int64      `db:"length_seconds" json:"length_seconds"`
	ViewCount     int64      `db:"view_count" json:"view_count"`
	AverageRating float64    `db:"average_rating" json:"average_rating"`
	ThumbnailURL  string     `db:"thumbnail_url" json:"thumbnail_url"`
	Category      string     `db:"category" json:"category"`
	IsPrivate     bool       `db:"is_private" json:"is_private"`
	PublishDate   *time.Time `db:"publish_date" json:"publish_date"`
	UploadDate    *time.Time `db:"upload_date" json:"upload_date"`
	LiveStartedAt *time.Time `db:"live_started_at" json:"live_started_at"`
	LiveEndedAt   *time.Time `db:"live_ended_at" json:"live_ended_at"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type VideoSlice []Video

func (v *VideoSlice) Scan(val interface{}) error {
	return scanJSON(v, val)
}
