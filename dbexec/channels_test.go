package dbexec

import (
	"os"
	"testing"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestChannelsInsertMany(t *testing.T) {
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	ex := ChannelsExecutor{db.MustBegin()}
	_, err = ex.InsertMany([]models.Channel{
		{
			ChannelID: "XXX",
			Name:      "YYY",
			ImageURL:  "https://placehold.jp/64x64.png?text=dummy",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestChannelsFind(t *testing.T) {
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	ex := ChannelsExecutor{db.MustBegin()}
	channel, err := ex.Find("CHANNEL_A")
	if err != nil {
		t.Fatal(err)
	}

	if channel.Name != "A channel" {
		t.Fatal("the name of the retrieved channel does not match")
	}
}
