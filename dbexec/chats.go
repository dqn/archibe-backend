package dbexec

import (
	"database/sql"
	"encoding/json"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
)

type ChatsExecutor struct {
	db *sqlx.DB
}

func (e *ChatsExecutor) InsertMany(chats []models.Chat) (sql.Result, error) {
	sql := `
	INSERT INTO chats (
		author_channel_id,
		video_id,
		timestamp,
		timestamp_usec,
		message_elements,
		purchase_amount,
		currency_unit,
		super_chat_context,
		created_at,
		updated_at
	)
	SELECT DISTINCT
		author_channel_id,
		video_id,
		timestamp,
		timestamp_usec,
		message_elements,
		COALESCE(purchase_amount, 0),
		COALESCE(currency_unit, ''),
		COALESCE(super_chat_context, '{}'),
		COALESCE(created_at, NOW()),
		COALESCE(updated_at, NOW())
	FROM
		jsonb_to_recordset($1) AS x(
			author_channel_id TEXT,
			video_id TEXT,
			timestamp TEXT,
			timestamp_usec TEXT,
			message_elements JSONB,
			purchase_amount NUMERIC,
			currency_unit TEXT,
			super_chat_context JSONB,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	ON CONFLICT(author_channel_id, video_id, timestamp_usec) DO NOTHING
	`

	b, err := json.Marshal(chats)
	if err != nil {
		return nil, err
	}

	return e.db.Exec(sql, string(b))
}
