package dbexec

import (
	"os"
	"testing"

	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestChatsInsertMany(t *testing.T) {
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	ex := ChatsExecutor{db}
	_, err = ex.InsertMany([]models.Chat{
		{
			AuthorChannelID: "AAAAA",
			VideoID:         "BBBBB",
			Timestamp:       "00:00",
			TimestampUsec:   "123456789",
			MessageElements: []models.MessageElement{
				{Type: "text", Text: "foo"},
				{Type: "text", Text: "bar"},
			},
			PurchaseAmount: 1000,
			CurrencyUnit:   "¥",
			SuperChatContext: &models.SuperChatContext{
				HeaderBackgroundColor: "ffffff",
				HeaderTextColor:       "ffffff",
				BodyBackgroundColor:   "ffffff",
				BodyTextColor:         "ffffff",
				AuthorNameTextColor:   "ffffff",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
