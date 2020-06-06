package models

import (
	"time"
)

type Chat struct {
	ID               int64               `db:"id" json:"id"`
	ChatID           string              `db:"chat_id" json:"chat_id"`
	AuthorChannelID  string              `db:"author_channel_id" json:"author_channel_id"`
	VideoID          string              `db:"video_id" json:"video_id"`
	Type             string              `db:"type" json:"type"`
	Timestamp        string              `db:"timestamp" json:"timestamp"`
	TimestampUsec    string              `db:"timestamp_usec" json:"timestamp_usec"`
	MessageElements  MessageElementSlice `db:"message_elements" json:"message_elements"`
	PurchaseAmount   float64             `db:"purchase_amount" json:"purchase_amount,omitempty"`
	CurrencyUnit     string              `db:"currency_unit" json:"currency_unit,omitempty"`
	SuperChatContext *SuperChatContext   `db:"super_chat_context" json:"super_chat_context,omitempty"`
	CreatedAt        *time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt        *time.Time          `db:"updated_at" json:"updated_at"`
	Channel          *Channel            `db:"channel" json:"channel,omitempty"`
	Badges           BadgeSlice          `db:"badges" json:"badges,omitempty"`
}

type ChatSlice []Chat

type MessageElement struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	Label    string `json:"label,omitempty"`
}

type MessageElementSlice []MessageElement

func (m *MessageElementSlice) Scan(val interface{}) error {
	return scanJSON(m, val)
}

type SuperChatContext struct {
	HeaderBackgroundColor string `json:"header_background_color,omitempty"`
	HeaderTextColor       string `json:"header_text_color,omitempty"`
	BodyBackgroundColor   string `json:"body_background_color,omitempty"`
	BodyTextColor         string `json:"body_text_color,omitempty"`
	AuthorNameTextColor   string `json:"author_name_text_color,omitempty"`
}

func (s *SuperChatContext) Scan(val interface{}) error {
	return scanJSON(s, val)
}
