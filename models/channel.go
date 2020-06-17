package models

import "time"

type Channel struct {
	ID                 int64                         `db:"id" json:"id"`
	ChannelID          string                        `db:"channel_id" json:"channel_id"`
	Name               string                        `db:"name" json:"name"`
	ImageURL           string                        `db:"image_url" json:"image_url"`
	CreatedAt          *time.Time                    `db:"created_at" json:"created_at"`
	UpdatedAt          *time.Time                    `db:"updated_at" json:"updated_at"`
	SentChatCount      int64                         `db:"sent_chat_count" json:"sent_chat_count"`
	ReceivedChatCount  int64                         `db:"received_chat_count" json:"received_chat_count"`
	SentSuperChats     SuperChatPerCurrencyUnitSlice `db:"sent_super_chats" json:"sent_super_chats"`
	ReceivedSuperChats SuperChatPerCurrencyUnitSlice `db:"received_super_chats" json:"received_super_chats"`
	Badges             BadgeSlice                    `db:"badges" json:"badges"`
	Videos             VideoSlice                    `db:"videos" json:"videos"`
}

type ChannelSlice []Channel
