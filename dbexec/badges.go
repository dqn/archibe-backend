package dbexec

import (
	"database/sql"
	"encoding/json"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
)

type BadgesExecutor struct {
	db *sqlx.DB
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
		image_url = EXCLUDED.image_url,
		label = EXCLUDED.label,
		updated_at = EXCLUDED.updated_at
	`

	b, err := json.Marshal(badges)
	if err != nil {
		return nil, err
	}

	return e.db.Exec(sql, string(b))
}

func (e *BadgesExecutor) FindByChannelID(channelID string) ([]models.Badge, error) {
	sql := `
	SELECT
		owner_channel_id,
		liver_channel_id,
		badge_type,
		image_url,
		label,
		created_at,
		updated_at
	FROM
		badges
	WHERE
		owner_channel_id = $1
	`

	badges := []models.Badge{}
	if err := e.db.Select(&badges, sql, channelID); err != nil {
		return nil, err
	}

	return badges, nil
}
