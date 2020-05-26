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
}

func NewExecutor(db *sqlx.DB) *DBExecutor {
	return &DBExecutor{
		DB:       db,
		Channels: &ChannelsExecutor{db},
		Videos:   &VideosExecutor{db},
		Chats:    &ChatsExecutor{db},
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
	;`

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
		channel_id,
		video_id,
		created_at,
		updated_at
	) VALUES (
		:channel_id,
		:video_id,
		COALESCE(:created_at, NOW()),
		COALESCE(:updated_at, NOW())
	)`

	return e.DB.Exec(sql, video)
}

type ChatsExecutor struct {
	DB *sqlx.DB
}

func (e *ChatsExecutor) InsertMany(chats []models.Chat) (sql.Result, error) {
	sql := `
	INSERT INTO chats (
		channel_id,
		video_id,
		timestamp,
		timestamp_usec,
		message_elements,
		purchase_amount,
		currency_unit,
		is_moderator,
		badge,
		created_at,
		updated_at
	)
	SELECT
		channel_id,
		video_id,
		timestamp,
		timestamp_usec,
		message_elements,
		COALESCE(purchase_amount, DEFAULT),
		COALESCE(currency_unit, DEFAULT),
		COALESCE(is_moderator, DEFAULT),
		COALESCE(badge, DEFAULT),
		COALESCE(created_at, NOW()),
		COALESCE(updated_at, NOW())
	FROM
    jsonb_to_recordset($1) AS x(
			channel_id TEXT,
			video_id TEXT,
			timestamp TEXT,
			timestamp_usec TEXT,
			message_elements JSONB,
			purchase_amount NUMERIC,
			currency_unit TEXT,
			is_moderator BOOLEAN,
			badge JSONB,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	;`

	b, err := json.Marshal(chats)
	if err != nil {
		return nil, err
	}

	return e.DB.Exec(sql, string(b))
}
