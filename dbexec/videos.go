package dbexec

import (
	"database/sql"
	"strings"

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
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		view_count = EXCLUDED.view_count,
		average_rating = EXCLUDED.average_rating,
		thumbnail_url = EXCLUDED.thumbnail_url,
		category = EXCLUDED.category,
		is_private = EXCLUDED.is_private,
		publish_date = EXCLUDED.publish_date,
		updated_at = EXCLUDED.updated_at
	`

	return e.db.NamedExec(sql, video)
}

type VideosQuery struct {
	Q       string
	Channel string
	Order   string
	Limit   uint64
	Offset  uint64
}

func (e *VideosExecutor) FindByQuery(query *VideosQuery) ([]models.Video, error) {
	order := strings.ToUpper(query.Order)
	if order != "DESC" {
		order = "ASC"
	}

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
		t1.updated_at,
		t2.id AS "channel.id",
		t2.channel_id AS "channel.channel_id",
		t2.name AS "channel.name",
		t2.image_url AS "channel.image_url",
		t2.created_at AS "channel.created_at",
		t2.updated_at AS "channel.updated_at"
	FROM
		videos AS t1
		INNER JOIN channels AS t2 ON (
			t1.channel_id = t2.channel_id
		)
	WHERE
		(
			$1 = ''
			OR t1.title ~ $1
		) AND (
			$2 = ''
			OR t1.channel_id = $2
		)
	ORDER BY
		t1.publish_date ` + order + `
	LIMIT
		$3
	OFFSET
		$4
	`

	videos := []models.Video{}
	if err := e.db.Select(&videos, sql, query.Q, query.Channel, query.Limit, query.Offset); err != nil {
		return nil, err
	}

	return videos, nil
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
					v1.video_id = $1
					AND v1.currency_unit != ''
				GROUP BY
					v1.currency_unit
			) AS u1
		) AS received_super_chats,
		t2.id AS "channel.id",
		t2.channel_id AS "channel.channel_id",
		t2.name AS "channel.name",
		t2.image_url AS "channel.image_url",
		t2.created_at AS "channel.created_at",
		t2.updated_at AS "channel.updated_at"
	FROM
		videos AS t1
		INNER JOIN channels AS t2 ON (
			t1.channel_id = t2.channel_id
		)
	WHERE
		video_id = $1
	`

	var video models.Video
	if err := e.db.Get(&video, sql, videoID); err != nil {
		return nil, err
	}

	return &video, nil
}
