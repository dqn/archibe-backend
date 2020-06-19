package dbexec

import (
	"os"
	"testing"
	"time"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestVideosInsertOne(t *testing.T) {
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	ex := VideosExecutor{db.MustBegin()}

	now := time.Now()
	_, err = ex.InsertOne(&models.Video{
		VideoID:       "AAA",
		ChannelID:     "BBB",
		Title:         "CCC",
		Description:   "DDD",
		LengthSeconds: 42,
		ViewCount:     42,
		AverageRating: 3.14,
		ThumbnailURL:  "EEE",
		Category:      "FFF",
		IsPrivate:     false,
		PublishDate:   &now,
		UploadDate:    &now,
		LiveStartedAt: &now,
		LiveEndedAt:   &now,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVideosFind(t *testing.T) {
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	ex := VideosExecutor{db.MustBegin()}

	_, err = ex.Find("VIDEO_A")
	if err != nil {
		t.Fatal(err)
	}
}
