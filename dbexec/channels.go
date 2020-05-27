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

	return e.db.Exec(sql, string(b))
}
