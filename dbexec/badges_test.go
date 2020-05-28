package dbexec

import (
	"fmt"
	"os"
	"testing"

	"github.com/dqn/tubekids/models"
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
			OwnerChannelID: "XXX",
			LiverChannelID: "YYY",
			BadgeType:      "member",
			ImageURL:       "https://placehold.jp/64x64.png?text=dummy",
			Label:          ":dummy:",
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

	fmt.Printf("%#v\n", badges)
}
