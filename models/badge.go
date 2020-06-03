package models

import "time"

type Badge struct {
	ID        int64      `db:"id" json:"id"`
	ChatID    string     `db:"chat_id" json:"chat_id"`
	BadgeType string     `db:"badge_type" json:"badge_type"`
	ImageURL  string     `db:"image_url" json:"image_url"`
	Label     string     `db:"label" json:"label"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

type BadgeSlice []Badge

func (m *BadgeSlice) Scan(val interface{}) error {
	return scanJSON(m, val)
}
