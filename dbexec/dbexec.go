package dbexec

import (
	"database/sql"
	"encoding/json"

	"github.com/dqn/tubekids/models"
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

type ChannelsExecutor struct {
	DB *sqlx.DB
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
	SELECT DISTINCT
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
		channel_id = EXCLUDED.channel_id,
		name = EXCLUDED.name,
		image_url = EXCLUDED.image_url,
		updated_at = EXCLUDED.updated_at
	`

	b, err := json.Marshal(channels)
	if err != nil {
		return nil, err
	}

	return e.DB.Exec(sql, string(b))
}

type VideosExecutor struct {
	DB *sqlx.DB
}

func (e *VideosExecutor) InsertOne(video *models.Video) (sql.Result, error) {
	sql := `
	INSERT INTO videos (
		video_id,
		channel_id,
		created_at,
		updated_at
	) VALUES (
		:video_id,
		:channel_id,
		COALESCE(:created_at, NOW()),
		COALESCE(:updated_at, NOW())
	)
	ON CONFLICT(video_id) DO NOTHING
	`

	return e.DB.NamedExec(sql, video)
}

type ChatsExecutor struct {
	DB *sqlx.DB
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
		created_at,
		updated_at
	)
	SELECT DISTINCT
		author_channel_id,
		video_id,
		timestamp,
		timestamp_usec,
		message_elements,
		purchase_amount,
		currency_unit,
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
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	ON CONFLICT(author_channel_id, video_id, timestamp_usec) DO NOTHING
	`

	b, err := json.Marshal(chats)
	if err != nil {
		return nil, err
	}

	return e.DB.Exec(sql, string(b))
}

type BadgesExecutor struct {
	DB *sqlx.DB
}

func (e *BadgesExecutor) InsertMany(badges []models.Badge) (sql.Result, error) {
	sql := `
	INSERT INTO badges (
		owner_channel_id,
		liver_channel_id,
		badge_type,
		image_url,
		label,
		created_at,
		updated_at
	)
	SELECT DISTINCT
		owner_channel_id,
		liver_channel_id,
		badge_type,
		image_url,
		label,
		COALESCE(created_at, NOW()),
		COALESCE(updated_at, NOW())
	FROM
    jsonb_to_recordset($1) AS x(
			owner_channel_id TEXT,
			liver_channel_id TEXT,
			badge_type TEXT,
			image_url TEXT,
			label TEXT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	ON CONFLICT (owner_channel_id, liver_channel_id, badge_type) DO UPDATE SET
		owner_channel_id = EXCLUDED.owner_channel_id,
		liver_channel_id = EXCLUDED.liver_channel_id,
		badge_type = EXCLUDED.badge_type,
		image_url = EXCLUDED.image_url,
		label = EXCLUDED.label,
		updated_at = EXCLUDED.updated_at
	`

	b, err := json.Marshal(badges)
	if err != nil {
		return nil, err
	}

	return e.DB.Exec(sql, string(b))
}
