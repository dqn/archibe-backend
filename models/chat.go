package models

import "time"

type Chat struct {
	ID              int64            `db:"id"`
	ChannelID       string           `db:"channel_id"`
	VideoID         string           `db:"video_id"`
	Timestamp       string           `db:"timestamp"`
	TimestampUsec   string           `db:"timestamp_usec"`
	MessageElements []MessageElement `db:"message_elements"`
	PurchaseAmount  float64          `db:"purchase_amount"`
	CurrencyUnit    string           `db:"currency_unit"`
	IsModerator     bool             `db:"is_moderator"`
	Badge           Badge            `db:"badge"`
	CreatedAt       time.Time        `db:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at"`
}

type MessageElement struct {
	Type  string
	Text  string
	Label string
	URL   string
}

type Badge struct {
	Label string
	URL   string
}
