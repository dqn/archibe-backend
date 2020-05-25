package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dqn/chatlog"
	"github.com/dqn/chatlog/chat"
	_ "github.com/lib/pq"
)

type Channel struct {
	ChannelID string
	Name      string
	imageURL  string
}

type Chat struct {
	ChannelID       string
	Timestamp       string
	TimestampUsec   string
	MessageElements []MessageElement
	PurchaseAmount  float64
	CurrencyUnit    string
	IsModerator     bool
	Badge           Badge
}

type MessageElement struct {
	Type  string
	Text  string
	Label string
	URL   string
}

type Badge struct {
	Label string
	URL   string
}

func parseNagesen(str string) (string, float64, error) {
	unit := strings.TrimRight(str, "0123456789.,")
	s := strings.TrimLeft(str, unit)
	s = strings.ReplaceAll(s, ",", "")
	amount, err := strconv.ParseFloat(s, 64)
	unit = strings.ReplaceAll(unit, "￥", "¥")
	return unit, amount, err
}

func parseMessage(message *chat.Message) ([]MessageElement, error) {
	me := make([]MessageElement, 0, len(message.Runs))
	for _, v := range message.Runs {
		var m MessageElement
		switch {
		case v.Emoji.EmojiId != "":
			m.Type = "emoji"
			m.Label = v.Emoji.Image.Accessibility.AccessibilityData.Label
			m.URL = v.Emoji.Image.Thumbnails[1].URL
		case v.Text != "":
			m.Type = "text"
			m.Text = v.Text
		default:
			err := fmt.Errorf("unknown message: %#v", v)
			return nil, err
		}

		me = append(me, m)
	}

	return me, nil
}

func parseAuthorBadges(badges []chat.AuthorBadge) (bool, *Badge) {
	var (
		isModerator bool
		badge       Badge
	)
	for _, b := range badges {
		if b.LiveChatAuthorBadgeRenderer.Icon.IconType == "MODERATOR" {
			isModerator = true
			continue
		}
		badge = Badge{
			Label: b.LiveChatAuthorBadgeRenderer.Accessibility.AccessibilityData.Label,
			URL:   b.LiveChatAuthorBadgeRenderer.CustomThumbnail.Thumbnails[1].URL,
		}
	}

	return isModerator, &badge
}

func run() error {
	if len(os.Args) != 3 {
		os.Exit(1)
	}

	dsn := os.Args[1]
	videoID := os.Args[2]

	cl, err := chatlog.New(videoID)
	if err != nil {
		return err
	}

	channels := make(map[string]Channel, 1024)
	chats := make([]Chat, 1024)

	for cl.Continuation != "" {
		cas, err := cl.Fecth()
		if err != nil {
			return err
		}

		for _, ca := range cas {
			for _, rcia := range ca.ReplayChatItemAction.Actions {
				item := rcia.AddChatItemAction.Item
				switch {
				case item.LiveChatTextMessageRenderer.ID != "":
					renderer := item.LiveChatTextMessageRenderer

					channels[renderer.AuthorExternalChannelId] = Channel{
						ChannelID: renderer.AuthorExternalChannelId,
						Name:      renderer.AuthorName.SimpleText,
						imageURL:  renderer.AuthorPhoto.Thumbnails[1].URL,
					}

					me, err := parseMessage(&renderer.Message)
					if err != nil {
						return err
					}

					isModerator, badge := parseAuthorBadges(renderer.AuthorBadges)

					chats = append(chats, Chat{
						ChannelID:       renderer.AuthorExternalChannelId,
						Timestamp:       renderer.TimestampText.SimpleText,
						TimestampUsec:   renderer.TimestampUsec,
						MessageElements: me,
						IsModerator:     isModerator,
						Badge:           *badge,
					})

				case item.LiveChatPaidMessageRenderer.ID != "":
					renderer := item.LiveChatPaidMessageRenderer

					channels[renderer.AuthorExternalChannelId] = Channel{
						ChannelID: renderer.AuthorExternalChannelId,
						Name:      renderer.AuthorName.SimpleText,
						imageURL:  renderer.AuthorPhoto.Thumbnails[1].URL,
					}

					me, err := parseMessage(&renderer.Message)
					if err != nil {
						return err
					}

					unit, amount, err := parseNagesen(renderer.PurchaseAmountText.SimpleText)
					if err != nil {
						return err
					}

					chats = append(chats, Chat{
						ChannelID:       renderer.AuthorExternalChannelId,
						Timestamp:       renderer.TimestampText.SimpleText,
						TimestampUsec:   renderer.TimestampUsec,
						MessageElements: me,
						PurchaseAmount:  amount,
						CurrencyUnit:    unit,
					})
				}
			}
		}
	}

	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	rows, err := pool.Query("SELECT 1;")
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return err
		}
		fmt.Println(id)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
