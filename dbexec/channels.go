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
	SELECT
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
		(SELECT COUNT(*) FROM chats WHERE chats.author_channel_id = $1) AS sent_chat_count,
		(SELECT COUNT(*) FROM chats INNER JOIN videos ON chats.video_id = videos.video_id WHERE videos.channel_id = $1) AS received_chat_count,
		t1.created_at,
		t1.updated_at,
		jsonb_agg(DISTINCT jsonb_build_object(
			'badge_type',
			t2.badge_type,
			'image_url',
			t2.image_url,
			'label',
			t2.label
		)) AS badges,
		jsonb_agg(DISTINCT jsonb_build_object(
			'video_id',
			t3.video_id
		)) AS videos
	FROM
		channels AS t1
		INNER JOIN badges AS t2 ON (
			t1.channel_id = t2.owner_channel_id
		)
		INNER JOIN videos AS t3 ON (
			t1.channel_id = t3.channel_id
		)
	WHERE
		t1.channel_id = $1
	GROUP BY
		t1.id
	`

	var channel models.Channel
	err := e.db.Get(&channel, sql, channelID)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}
