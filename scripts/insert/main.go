package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dqn/chatlog"
	"github.com/dqn/chatlog/chat"
	"github.com/dqn/tubekids/dbexec"
	"github.com/dqn/tubekids/models"
	"github.com/dqn/ytvi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func toHex(d int) string {
	return strconv.FormatInt(int64(d), 16)
}

func retrieveImageURL(thumbnails []chat.Thumbnail) string {
	return thumbnails[len(thumbnails)-1].URL
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
			m.ImageURL = retrieveImageURL(v.Emoji.Image.Thumbnails)
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
	if len(os.Args) != 3 {
		return fmt.Errorf("invalid arguments")
	}

	dsn := os.Args[1]
	videoID := os.Args[2]

	cl, err := chatlog.New(videoID)
	if err != nil {
		return err
	}

	bufsize := 1024
	channels := make([]models.Channel, 0, bufsize)
	chats := make([]models.Chat, 0, bufsize)
	badges := make([]models.Badge, 0, bufsize)

	appendedChannels := make(map[string]struct{}, bufsize)

	fmt.Println("start fetching chats...")

	err = processEachChatItem(cl, func(item *chat.ChatItem) error {
		switch {
		case item.LiveChatTextMessageRenderer.ID != "":
			renderer := item.LiveChatTextMessageRenderer

			if _, ok := appendedChannels[renderer.AuthorExternalChannelID]; !ok {
				channels = append(channels, models.Channel{
					ChannelID: renderer.AuthorExternalChannelID,
					Name:      renderer.AuthorName.SimpleText,
					ImageURL:  retrieveImageURL(renderer.AuthorPhoto.Thumbnails),
				})
				appendedChannels[renderer.AuthorExternalChannelID] = struct{}{}
			}

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
						ChatID:    item.LiveChatTextMessageRenderer.ID,
						BadgeType: "moderator",
					})

				default:
					badges = append(badges, models.Badge{
						ChatID:    item.LiveChatTextMessageRenderer.ID,
						BadgeType: "member",
						ImageURL:  retrieveImageURL(b.LiveChatAuthorBadgeRenderer.CustomThumbnail.Thumbnails),
						Label:     b.LiveChatAuthorBadgeRenderer.Accessibility.AccessibilityData.Label,
					})
				}
			}

			chats = append(chats, models.Chat{
				ChatID:          item.LiveChatTextMessageRenderer.ID,
				AuthorChannelID: renderer.AuthorExternalChannelID,
				VideoID:         videoID,
				Type:            "chat",
				Timestamp:       renderer.TimestampText.SimpleText,
				TimestampUsec:   renderer.TimestampUsec,
				MessageElements: me,
			})

		case item.LiveChatPaidMessageRenderer.ID != "":
			renderer := item.LiveChatPaidMessageRenderer

			if _, ok := appendedChannels[renderer.AuthorExternalChannelID]; !ok {
				channels = append(channels, models.Channel{
					ChannelID: renderer.AuthorExternalChannelID,
					Name:      renderer.AuthorName.SimpleText,
					ImageURL:  retrieveImageURL(renderer.AuthorPhoto.Thumbnails),
				})
				appendedChannels[renderer.AuthorExternalChannelID] = struct{}{}
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
				ChatID:          item.LiveChatPaidMessageRenderer.ID,
				AuthorChannelID: renderer.AuthorExternalChannelID,
				VideoID:         videoID,
				Type:            "super_chat",
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

	pr, err := ytvi.GetVideoInfo(videoID)
	if err != nil {
		return err
	}

	pmr := pr.Microformat.PlayerMicroformatRenderer

	ownerChannel := models.Channel{
		ChannelID: pmr.ExternalChannelID,
		Name:      pmr.OwnerChannelName,
		ImageURL:  pmr.OwnerProfileURL,
	}

	channels = append(channels, ownerChannel)

	lengthSeconds, err := strconv.ParseInt(pmr.LengthSeconds, 10, 64)
	if err != nil {
		return err
	}
	viewCount, err := strconv.ParseInt(pmr.ViewCount, 10, 64)
	if err != nil {
		return err
	}
	publishDate, err := time.Parse("2006-01-02", pmr.PublishDate)
	if err != nil {
		return err
	}
	uploadDate, err := time.Parse("2006-01-02", pmr.UploadDate)
	if err != nil {
		return err
	}

	thumbnails := pr.VideoDetails.Thumbnail.Thumbnails
	thumbnailURL := thumbnails[len(thumbnails)-1].URL

	video := models.Video{
		VideoID:       videoID,
		ChannelID:     pmr.ExternalChannelID,
		Title:         pmr.Title.SimpleText,
		Description:   pmr.Description.SimpleText,
		LengthSeconds: lengthSeconds,
		ViewCount:     viewCount,
		AverageRating: pr.VideoDetails.AverageRating,
		ThumbnailURL:  thumbnailURL,
		Category:      pmr.Category,
		IsPrivate:     pr.VideoDetails.IsPrivate,
		PublishDate:   &publishDate,
		UploadDate:    &uploadDate,
		LiveStartedAt: &pmr.LiveBroadcastDetails.StartTimestamp,
		LiveEndedAt:   &pmr.LiveBroadcastDetails.EndTimestamp,
	}

	_, err = dbx.Videos.InsertOne(&video)
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
