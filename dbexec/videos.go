package dbexec

import (
	"database/sql"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
)

type VideosExecutor struct {
	db *sqlx.DB
}

func (e *VideosExecutor) InsertOne(video *models.Video) (sql.Result, error) {
	sql := `
	INSERT INTO videos (
		video_id,
		channel_id,
		title,
		description,
		length_seconds,
		view_count,
		average_rating,
		thumbnail_url,
		category,
		is_private,
		publish_date,
		upload_date,
		live_started_at,
		live_ended_at,
		created_at,
		updated_at
	) VALUES (
		:video_id,
		:channel_id,
		:title,
		:description,
		:length_seconds,
		:view_count,
		:average_rating,
		:thumbnail_url,
		:category,
		:is_private,
		:publish_date,
		:upload_date,
		:live_started_at,
		:live_ended_at,
		COALESCE(:created_at, NOW()),
		COALESCE(:updated_at, NOW())
	)
	ON CONFLICT(video_id) DO UPDATE SET
		channel_id = EXCLUDED.channel_id,
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		length_seconds = EXCLUDED.length_seconds,
		view_count = EXCLUDED.view_count,
		average_rating = EXCLUDED.average_rating,
		thumbnail_url = EXCLUDED.thumbnail_url,
		category = EXCLUDED.category,
		is_private = EXCLUDED.is_private,
		publish_date = EXCLUDED.publish_date,
		upload_date = EXCLUDED.upload_date,
		live_started_at = EXCLUDED.live_started_at,
		live_ended_at = EXCLUDED.live_ended_at,
		updated_at = EXCLUDED.updated_at
	`

	return e.db.NamedExec(sql, video)
}
