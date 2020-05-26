package models

import "time"

type Chat struct {
	ID              int64            `db:"id" json:"id"`
	ChannelID       string           `db:"channel_id" json:"channel_id"`
	VideoID         string           `db:"video_id" json:"video_id"`
	Timestamp       string           `db:"timestamp" json:"timestamp"`
	TimestampUsec   string           `db:"timestamp_usec" json:"timestamp_usec"`
	MessageElements []MessageElement `db:"message_elements" json:"message_elements"`
	PurchaseAmount  float64          `db:"purchase_amount" json:"purchase_amount"`
	CurrencyUnit    string           `db:"currency_unit" json:"currency_unit"`
	IsModerator     bool             `db:"is_moderator" json:"is_moderator"`
	Badge           *Badge           `db:"badge" json:"badge"`
	CreatedAt       *time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time       `db:"updated_at" json:"updated_at"`
}

type MessageElement struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Badge struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}
