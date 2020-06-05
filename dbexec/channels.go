package dbexec

import (
	"database/sql"
	"encoding/json"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
)

type ChannelsExecutor struct {
	db *sqlx.DB
}

func (e *ChannelsExecutor) InsertMany(channels []models.Channel) (sql.Result, error) {
	sql := `
	INSERT INTO channels (
		channel_id,
		name,
		image_url,
		created_at,
		updated_at
	)
	SELECT
		channel_id,
		name,
		image_url,
		COALESCE(created_at, NOW()),
		COALESCE(updated_at, NOW())
	FROM
		jsonb_to_recordset($1) AS x(
			channel_id TEXT,
			name TEXT,
			image_url TEXT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	ON CONFLICT (channel_id) DO UPDATE SET
		name = EXCLUDED.name,
		image_url = EXCLUDED.image_url,
		updated_at = EXCLUDED.updated_at
	`

	b, err := json.Marshal(channels)
	if err != nil {
		return nil, err
	}

	return e.db.Exec(sql, string(b))
}

func (e *ChannelsExecutor) Find(channelID string) (*models.Channel, error) {
	sql := `
	SELECT
		t1.id,
		t1.channel_id,
		t1.name,
		t1.image_url,
		(
			SELECT
				COUNT(*)
			FROM
				chats AS u1
			WHERE
				u1.author_channel_id = $1
		) AS sent_chat_count,
		(
			SELECT
				COUNT(*)
			FROM
				chats AS u1
				INNER JOIN videos AS u2 ON (
					u1.video_id = u2.video_id
				)
			WHERE
				u2.channel_id = $1
		) AS received_chat_count,
		t1.created_at,
		t1.updated_at,
		(
			SELECT DISTINCT
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
				INNER JOIN chats AS u2 ON (
					u1.chat_id = u2.chat_id
				)
			WHERE
				u2.author_channel_id = $1
		) AS badges,
		(
			SELECT DISTINCT
				COALESCE(jsonb_agg(jsonb_build_object(
					'video_id',
					u1.video_id
				) ORDER BY u1.created_at), '[]')
			FROM
				videos AS u1
			WHERE
				u1.channel_id = $1
		) AS videos
	FROM
		channels AS t1
	WHERE
		t1.channel_id = $1
	`

	var channel models.Channel
	err := e.db.Get(&channel, sql, channelID)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}
