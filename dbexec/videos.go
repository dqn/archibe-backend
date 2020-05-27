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

	return e.db.NamedExec(sql, video)
}
