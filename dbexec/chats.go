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
		chat_id,
		author_channel_id,
		video_id,
		type,
		timestamp,
		timestamp_usec,
		message_elements,
		purchase_amount,
		currency_unit,
		super_chat_context,
		created_at,
		updated_at
	)
	SELECT
		chat_id,
		author_channel_id,
		video_id,
		type,
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
			chat_id TEXT,
			author_channel_id TEXT,
			video_id TEXT,
			type TEXT,
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

type ChatsQuery struct {
	Q       string
	Channel string
	Video   string
	Limit   uint64
	Offset  uint64
}

func (e *ChatsExecutor) FindByQuery(query *ChatsQuery) ([]models.Chat, error) {
	sql := `
	SELECT
		t1.chat_id,
		t1.author_channel_id,
		t1.video_id,
		t1.type,
		t1.timestamp,
		t1.timestamp_usec,
		t1.message_elements,
		t1.purchase_amount,
		t1.currency_unit,
		t1.super_chat_context,
		t1.created_at,
		t1.updated_at,
		t2.channel_id AS "channel.channel_id",
		t2.name AS "channel.name",
		t2.image_url AS "channel.image_url",
		t2.created_at AS "channel.created_at",
		t2.updated_at AS "channel.updated_at",
		(
			SELECT
				COALESCE(jsonb_agg(jsonb_build_object(
					'badge_type',
					u1.badge_type,
					'image_url',
					u1.image_url,
					'label',
					u1.label
				) ORDER BY u1.created_at), '[]')
			FROM
				badges AS u1
			WHERE
				u1.chat_id = t1.chat_id
		) AS badges
	FROM
		chats AS t1
		INNER JOIN channels AS t2 ON (
			t1.author_channel_id = t2.channel_id
		)
		INNER JOIN videos AS t3 ON (
			t1.video_id = t3.video_id
		)
	WHERE
		(
			$1 = ''
			OR EXISTS (SELECT 1 FROM jsonb_to_recordset(t1.message_elements) as x(text TEXT) WHERE text IS NOT NULL AND text ~ $1)
		) AND (
			$2 = ''
			OR t1.author_channel_id = $2
		) AND (
			$3 = ''
			OR t1.video_id = $3
		)
	ORDER BY
		t1.timestamp_usec DESC
	LIMIT
		$4
	OFFSET
		$5
	`

	chats := []models.Chat{}
	if err := e.db.Select(&chats, sql, query.Q, query.Channel, query.Video, query.Limit, query.Offset); err != nil {
		return nil, err
	}

	return chats, nil
}
