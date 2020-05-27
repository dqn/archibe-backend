package dbexec

import (
	"github.com/jmoiron/sqlx"
)

type DBExecutor struct {
	DB       *sqlx.DB
	Channels *ChannelsExecutor
	Videos   *VideosExecutor
	Chats    *ChatsExecutor
	Badges   *BadgesExecutor
}

func NewExecutor(db *sqlx.DB) *DBExecutor {
	return &DBExecutor{
		DB:       db,
		Channels: &ChannelsExecutor{db},
		Videos:   &VideosExecutor{db},
		Chats:    &ChatsExecutor{db},
		Badges:   &BadgesExecutor{db},
	}
}
