package dbexec

import (
	"os"
	"testing"

	"github.com/dqn/archibe/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestBadgesInsertMany(t *testing.T) {
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	ex := BadgesExecutor{db}
	_, err = ex.InsertMany([]models.Badge{
		{
			ChatID:    "XXX",
			BadgeType: "member",
			ImageURL:  "https://placehold.jp/64x64.png?text=dummy",
			Label:     ":dummy:",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBadgesFindByChannelID(t *testing.T) {
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	ex := BadgesExecutor{db}
	badges, err := ex.FindByChannelID("CHANNEL_A")
	if err != nil {
		t.Fatal(err)
	}

	if len(badges) != 1 {
		t.Fatal("number of retrieved badges do not match")
	}
}
