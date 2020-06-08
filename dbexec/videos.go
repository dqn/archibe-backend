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

func (e *VideosExecutor) Find(videoID string) (*models.Video, error) {
	sql := `
	SELECT
		t1.id,
		t1.video_id,
		t1.channel_id,
		t1.title,
		t1.description,
		t1.length_seconds,
		t1.view_count,
		t1.average_rating,
		t1.thumbnail_url,
		t1.category,
		t1.is_private,
		t1.publish_date,
		t1.upload_date,
		t1.live_started_at,
		t1.live_ended_at,
		t1.created_at,
		t1.updated_at
	FROM
		videos AS t1
	WHERE
		video_id = $1
	`

	var video models.Video
	if err := e.db.Get(&video, sql, videoID); err != nil {
		return nil, err
	}

	return &video, nil
}
