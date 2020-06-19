package dbexec

import (
	"database/sql"
	"encoding/json"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
)

type BadgesExecutor struct {
	tx *sqlx.Tx
}

func (e *BadgesExecutor) InsertMany(badges []models.Badge) (sql.Result, error) {
	sql := `
	INSERT INTO badges (
		chat_id,
		badge_type,
		image_url,
		label,
		created_at,
		updated_at
	)
	SELECT
		chat_id,
		badge_type,
		image_url,
		label,
		COALESCE(created_at, NOW()),
		COALESCE(updated_at, NOW())
	FROM
		jsonb_to_recordset($1) AS x(
			chat_id TEXT,
			badge_type TEXT,
			image_url TEXT,
			label TEXT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	ON CONFLICT (chat_id, badge_type) DO UPDATE SET
		badge_type = EXCLUDED.badge_type,
		image_url = EXCLUDED.image_url,
		label = EXCLUDED.label,
		updated_at = EXCLUDED.updated_at
	`

	b, err := json.Marshal(badges)
	if err != nil {
		return nil, err
	}

	return e.tx.Exec(sql, string(b))
}

func (e *BadgesExecutor) FindByChannelID(channelID string) ([]models.Badge, error) {
	sql := `
	SELECT DISTINCT
		t1.badge_type,
		t1.image_url,
		t1.label
	FROM
		badges AS t1
		INNER JOIN chats AS t2 ON (
			t1.chat_id = t2.chat_id
		)
	WHERE
		t2.author_channel_id = $1
	`

	badges := []models.Badge{}
	if err := e.tx.Select(&badges, sql, channelID); err != nil {
		return nil, err
	}

	return badges, nil
}
