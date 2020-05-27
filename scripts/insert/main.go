package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dqn/chatlog"
	"github.com/dqn/chatlog/chat"
	"github.com/dqn/tubekids/dbexec"
	"github.com/dqn/tubekids/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func toHex(d int) string {
	return strconv.FormatInt(int64(d), 16)
}

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
			m.ImageURL = v.Emoji.Image.Thumbnails[1].URL
			m.Label = v.Emoji.Image.Accessibility.AccessibilityData.Label
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

func processEachChatItem(cl *chatlog.Chatlog, handler func(item *chat.ChatItem) error) error {
	for cl.Continuation != "" {
		continuationActions, err := cl.Fecth()
		if err != nil {
			return err
		}

		for _, ca := range continuationActions {
			for _, a := range ca.ReplayChatItemAction.Actions {
				if err = handler(&a.AddChatItemAction.Item); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func run() error {
	if len(os.Args) != 4 {
		os.Exit(1)
	}

	dsn := os.Args[1]
	channelID := os.Args[2]
	videoID := os.Args[3]

	cl, err := chatlog.New(videoID)
	if err != nil {
		return err
	}

	bufsize := 1024
	channels := make([]models.Channel, 0, bufsize)
	chats := make([]models.Chat, 0, bufsize)
	badges := make([]models.Badge, 0, bufsize)

	fmt.Println("start fetching chats...")

	err = processEachChatItem(cl, func(item *chat.ChatItem) error {
		switch {
		case item.LiveChatTextMessageRenderer.ID != "":
			renderer := item.LiveChatTextMessageRenderer

			channels = append(channels, models.Channel{
				ChannelID: renderer.AuthorExternalChannelID,
				Name:      renderer.AuthorName.SimpleText,
				ImageURL:  renderer.AuthorPhoto.Thumbnails[1].URL,
			})

			me, err := parseMessage(&renderer.Message)
			if err != nil {
				return err
			}

			for _, b := range renderer.AuthorBadges {
				switch b.LiveChatAuthorBadgeRenderer.Icon.IconType {
				case "OWNER":
					// do nothing

				case "MODERATOR":
					badges = append(badges, models.Badge{
						OwnerChannelID: renderer.AuthorExternalChannelID,
						LiverChannelID: channelID,
						BadgeType:      "moderator",
					})

				default:
					badges = append(badges, models.Badge{
						OwnerChannelID: renderer.AuthorExternalChannelID,
						LiverChannelID: channelID,
						BadgeType:      "member",
						ImageURL:       b.LiveChatAuthorBadgeRenderer.CustomThumbnail.Thumbnails[1].URL,
						Label:          b.LiveChatAuthorBadgeRenderer.Accessibility.AccessibilityData.Label,
					})
				}
			}

			chats = append(chats, models.Chat{
				AuthorChannelID: renderer.AuthorExternalChannelID,
				VideoID:         videoID,
				Timestamp:       renderer.TimestampText.SimpleText,
				TimestampUsec:   renderer.TimestampUsec,
				MessageElements: me,
			})

		case item.LiveChatPaidMessageRenderer.ID != "":
			renderer := item.LiveChatPaidMessageRenderer

			channels = append(channels, models.Channel{
				ChannelID: renderer.AuthorExternalChannelID,
				Name:      renderer.AuthorName.SimpleText,
				ImageURL:  renderer.AuthorPhoto.Thumbnails[1].URL,
			})

			me, err := parseMessage(&renderer.Message)
			if err != nil {
				return err
			}

			unit, amount, err := parseNagesen(renderer.PurchaseAmountText.SimpleText)
			if err != nil {
				return err
			}

			chats = append(chats, models.Chat{
				AuthorChannelID: renderer.AuthorExternalChannelID,
				VideoID:         videoID,
				Timestamp:       renderer.TimestampText.SimpleText,
				TimestampUsec:   renderer.TimestampUsec,
				MessageElements: me,
				PurchaseAmount:  amount,
				CurrencyUnit:    unit,
				SuperChatContext: &models.SuperChatContext{
					HeaderBackgroundColor: toHex(renderer.HeaderBackgroundColor),
					HeaderTextColor:       toHex(renderer.HeaderTextColor),
					BodyBackgroundColor:   toHex(renderer.BodyBackgroundColor),
					BodyTextColor:         toHex(renderer.BodyTextColor),
					AuthorNameTextColor:   toHex(renderer.AuthorNameTextColor),
				},
			})
		}

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("start inserting to database...")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	dbx := dbexec.NewExecutor(db)

	_, err = dbx.Videos.InsertOne(&models.Video{
		VideoID:   videoID,
		ChannelID: channelID,
	})
	if err != nil {
		return err
	}

	_, err = dbx.Channels.InsertMany(channels)
	if err != nil {
		return err
	}

	_, err = dbx.Chats.InsertMany(chats)
	if err != nil {
		return err
	}

	_, err = dbx.Badges.InsertMany(badges)
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
