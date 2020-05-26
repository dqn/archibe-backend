package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dqn/chatlog"
	"github.com/dqn/chatlog/chat"
	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func parseNagesen(str string) (string, float64, error) {
	unit := strings.TrimRight(str, "0123456789.,")
	s := strings.TrimLeft(str, unit)
	s = strings.ReplaceAll(s, ",", "")
	amount, err := strconv.ParseFloat(s, 64)
	unit = strings.ReplaceAll(unit, "￥", "¥")
	return unit, amount, err
}

func parseMessage(message *chat.Message) ([]models.MessageElement, error) {
	me := make([]models.MessageElement, 0, len(message.Runs))
	for _, v := range message.Runs {
		var m models.MessageElement
		switch {
		case v.Emoji.EmojiID != "":
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

func parseAuthorBadges(badges []chat.AuthorBadge) (bool, *models.Badge) {
	var (
		isModerator bool
		badge       models.Badge
	)
	for _, b := range badges {
		if b.LiveChatAuthorBadgeRenderer.Icon.IconType == "MODERATOR" {
			isModerator = true
			continue
		}
		badge = models.Badge{
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

	channelsMemo := make(map[string]struct{}, 1024)
	channels := make([]models.Channel, 1024)
	chats := make([]models.Chat, 1024)

	fmt.Println("start fetching chats...")

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

					if _, ok := channelsMemo[renderer.AuthorExternalChannelID]; !ok {
						channels = append(channels, models.Channel{
							ChannelID: renderer.AuthorExternalChannelID,
							Name:      renderer.AuthorName.SimpleText,
							ImageURL:  renderer.AuthorPhoto.Thumbnails[1].URL,
						})
						channelsMemo[renderer.AuthorExternalChannelID] = struct{}{}
					}

					me, err := parseMessage(&renderer.Message)
					if err != nil {
						return err
					}

					isModerator, badge := parseAuthorBadges(renderer.AuthorBadges)

					chats = append(chats, models.Chat{
						ChannelID:       renderer.AuthorExternalChannelID,
						VideoID:         videoID,
						Timestamp:       renderer.TimestampText.SimpleText,
						TimestampUsec:   renderer.TimestampUsec,
						MessageElements: me,
						IsModerator:     isModerator,
						Badge:           *badge,
					})

				case item.LiveChatPaidMessageRenderer.ID != "":
					renderer := item.LiveChatPaidMessageRenderer

					if _, ok := channelsMemo[renderer.AuthorExternalChannelID]; !ok {
						channels = append(channels, models.Channel{
							ChannelID: renderer.AuthorExternalChannelID,
							Name:      renderer.AuthorName.SimpleText,
							ImageURL:  renderer.AuthorPhoto.Thumbnails[1].URL,
						})
						channelsMemo[renderer.AuthorExternalChannelID] = struct{}{}
					}

					me, err := parseMessage(&renderer.Message)
					if err != nil {
						return err
					}

					unit, amount, err := parseNagesen(renderer.PurchaseAmountText.SimpleText)
					if err != nil {
						return err
					}

					chats = append(chats, models.Chat{
						ChannelID:       renderer.AuthorExternalChannelID,
						VideoID:         videoID,
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

	fmt.Println("start inserting to Database...")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	_, err = db.NamedExec(
		`INSERT INTO videos (
			channel_id,
			video_id,
			created_at,
			updated_at
		) VALUES (
			:channel_id,
			:video_id,
			:created_at,
			:updated_at
		)`,
		models.Video{
			VideoID:   videoID,
			ChannelID: "TODO",
		},
	)
	if err != nil {
		return err
	}

	_, err = db.NamedExec(
		`INSERT INTO channels (
			channel_id,
			name,
			image_url,
			created_at,
			updated_at
		) VALUES (
			:channel_id,
			:name,
			:image_url,
			:created_at,
			:updated_at
		)`,
		channels,
	)
	if err != nil {
		return err
	}

	_, err = db.NamedExec(
		`INSERT INTO chats (
			channel_id,
			video_id,
			timestamp,
			timestamp_usec,
			-- message_elements,
			purchase_amount,
			currency_unit,
			is_moderator,
			-- badge,
			created_at,
			updated_at
		) VALUES (
			:channel_id,
			:video_id,
			:timestamp,
			:timestamp_usec,
			-- message_elements,
			COALESCE(:purchase_amount, DEFAULT),
			COALESCE(:currency_unit, DEFAULT),
			COALESCE(:is_moderator, DEFAULT),
			-- :badge,
			:created_at,
			:updated_at
		);`,
		chats,
	)
	if err != nil {
		return err
	}

	fmt.Println("completed!")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
