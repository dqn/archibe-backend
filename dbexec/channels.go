package dbexec

import (
	"database/sql"
	"encoding/json"

	"github.com/dqn/archibe/models"
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

type ChannelsQuery struct {
	Q      string
	Limit  uint64
	Offset uint64
}

func (e *ChannelsExecutor) FindByQuery(query *ChannelsQuery) ([]models.Channel, error) {
	sql := `
	SELECT
		t1.id,
		t1.channel_id,
		t1.name,
		t1.image_url,
		(
			SELECT
				COALESCE(jsonb_agg(jsonb_build_object(
					'badge_type',
					u1.badge_type,
					'image_url',
					u1.image_url,
					'label',
					u1.label
				) ORDER BY u1.channel_id), '[]')
			FROM
				(
					SELECT DISTINCT ON (v3.channel_id) -- select latest badge for each channel
						v1.badge_type,
						v1.image_url,
						v1.label,
						v3.channel_id
					FROM
						badges AS v1
						INNER JOIN chats AS v2 ON (
							v1.chat_id = v2.chat_id
						)
						INNER JOIN videos AS v3 ON (
							v2.video_id = v3.video_id
						)
					WHERE
						v1.badge_type != 'moderator'
						AND v2.author_channel_id = t1.channel_id
					ORDER BY
						v3.channel_id,
						v2.timestamp_usec DESC
				) AS u1
		) AS badges
	FROM
		channels AS t1
	WHERE
		$1 = ''
		OR t1.name ~ $1
	ORDER BY
		t1.channel_id
	LIMIT
		$2
	OFFSET
		$3
	`

	channels := []models.Channel{}
	if err := e.db.Select(&channels, sql, query.Q, query.Limit, query.Offset); err != nil {
		return nil, err
	}

	return channels, nil
}

func (e *ChannelsExecutor) Find(channelID string) (*models.Channel, error) {
	sql := `
	SELECT
		t1.id,
		t1.channel_id,
		t1.name,
		t1.image_url,
		t1.created_at,
		t1.updated_at,
		(
			SELECT
				COALESCE(jsonb_agg(u1.object), '[]')
			FROM (
				SELECT
					jsonb_build_object(
						'currency_unit',
						v1.currency_unit,
						'purchase_amount',
						SUM(v1.purchase_amount)
					) AS object
				FROM
					chats AS v1
				WHERE
					v1.author_channel_id = $1
					AND v1.currency_unit != ''
				GROUP BY
					v1.currency_unit
			) AS u1
		) AS sent_super_chats,
		(
			SELECT
				COALESCE(jsonb_agg(u1.object), '[]')
			FROM (
				SELECT
					jsonb_build_object(
						'currency_unit',
						v1.currency_unit,
						'purchase_amount',
						SUM(v1.purchase_amount)
					) AS object
				FROM
					chats AS v1
					INNER JOIN videos AS v2 ON (
						v1.video_id = v2.video_id
					)
				WHERE
					v2.channel_id = $1
					AND v1.currency_unit != ''
				GROUP BY
					v1.currency_unit
			) AS u1
		) AS received_super_chats,
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
		(
			SELECT
				COALESCE(jsonb_agg(jsonb_build_object(
					'badge_type',
					u1.badge_type,
					'image_url',
					u1.image_url,
					'label',
					u1.label
				) ORDER BY u1.channel_id), '[]')
			FROM
				(
					SELECT DISTINCT ON (v3.channel_id) -- select latest badge for each channel
						v1.badge_type,
						v1.image_url,
						v1.label,
						v3.channel_id
					FROM
						badges AS v1
						INNER JOIN chats AS v2 ON (
							v1.chat_id = v2.chat_id
						)
						INNER JOIN videos AS v3 ON (
							v2.video_id = v3.video_id
						)
					WHERE
						v1.badge_type != 'moderator'
						AND v2.author_channel_id = $1
					ORDER BY
						v3.channel_id,
						v2.timestamp_usec DESC
				) AS u1
		) AS badges
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
