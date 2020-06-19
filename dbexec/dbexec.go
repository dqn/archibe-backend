package dbexec

import (
	"github.com/jmoiron/sqlx"
)

type DBExecutor struct {
	Tx       *sqlx.Tx
	Channels *ChannelsExecutor
	Videos   *VideosExecutor
	Chats    *ChatsExecutor
	Badges   *BadgesExecutor
}

func NewExecutor(tx *sqlx.Tx) *DBExecutor {
	return &DBExecutor{
		Tx:       tx,
		Channels: &ChannelsExecutor{tx},
		Videos:   &VideosExecutor{tx},
		Chats:    &ChatsExecutor{tx},
		Badges:   &BadgesExecutor{tx},
	}
}
