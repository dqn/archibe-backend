package archive

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dqn/chatlog"
	"github.com/dqn/tubekids/lib"
	"github.com/dqn/tubekids/models"
)

type ArchiveFetcher struct {
	videoID     string
	result      *ArchiveResult
	channelMemo *lib.Memo
}

type ArchiveResult struct {
	Channels []models.Channel
	Chats    []models.Chat
	Badges   []models.Badge
}

func NewFetcher(videoID string) *ArchiveFetcher {
	return &ArchiveFetcher{videoID: videoID}
}

func (a *ArchiveFetcher) Fetch() (*ArchiveResult, error) {
	c := chatlog.New(a.videoID)

	var size uint64 = 1024
	a.result = &ArchiveResult{
		Channels: make([]models.Channel, 0, size),
		Chats:    make([]models.Chat, 0, size),
		Badges:   make([]models.Badge, 0, size),
	}
	a.channelMemo = lib.NewMemo(size)

	err := c.HandleChatItem(func(item *chatlog.ChatItem) error {
		switch {
		case item.LiveChatTextMessageRenderer.ID != "":
			return a.handleTextMessage(&item.LiveChatTextMessageRenderer)
		case item.LiveChatPaidMessageRenderer.ID != "":
			return a.handlePaidMessage(&item.LiveChatPaidMessageRenderer)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return a.result, nil
}

func (a *ArchiveFetcher) handleTextMessage(renderer *chatlog.LiveChatTextMessageRenderer) error {
	if !a.channelMemo.Exists(renderer.AuthorExternalChannelID) {
		a.result.Channels = append(a.result.Channels, models.Channel{
			ChannelID: renderer.AuthorExternalChannelID,
			Name:      renderer.AuthorName.SimpleText,
			ImageURL:  retrieveImageURL(renderer.AuthorPhoto.Thumbnails),
		})
		a.channelMemo.Add(renderer.AuthorExternalChannelID)
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
			a.result.Badges = append(a.result.Badges, models.Badge{
				ChatID:    renderer.ID,
				BadgeType: "moderator",
			})

		default:
			a.result.Badges = append(a.result.Badges, models.Badge{
				ChatID:    renderer.ID,
				BadgeType: "member",
				ImageURL:  retrieveImageURL(b.LiveChatAuthorBadgeRenderer.CustomThumbnail.Thumbnails),
				Label:     b.LiveChatAuthorBadgeRenderer.Accessibility.AccessibilityData.Label,
			})
		}
	}

	timestampUsec, _ := strconv.ParseInt(renderer.TimestampUsec, 10, 64)

	a.result.Chats = append(a.result.Chats, models.Chat{
		ChatID:          renderer.ID,
		AuthorChannelID: renderer.AuthorExternalChannelID,
		VideoID:         a.videoID,
		Type:            "chat",
		Timestamp:       renderer.TimestampText.SimpleText,
		TimestampUsec:   timestampUsec,
		MessageElements: me,
	})

	return nil
}

func (a *ArchiveFetcher) handlePaidMessage(renderer *chatlog.LiveChatPaidMessageRenderer) error {
	if !a.channelMemo.Exists(renderer.AuthorExternalChannelID) {
		a.result.Channels = append(a.result.Channels, models.Channel{
			ChannelID: renderer.AuthorExternalChannelID,
			Name:      renderer.AuthorName.SimpleText,
			ImageURL:  retrieveImageURL(renderer.AuthorPhoto.Thumbnails),
		})
		a.channelMemo.Add(renderer.AuthorExternalChannelID)
	}

	me, err := parseMessage(&renderer.Message)
	if err != nil {
		return err
	}

	unit, amount, err := parseSuperChat(renderer.PurchaseAmountText.SimpleText)
	if err != nil {
		return err
	}

	timestampUsec, _ := strconv.ParseInt(renderer.TimestampUsec, 10, 64)

	a.result.Chats = append(a.result.Chats, models.Chat{
		ChatID:          renderer.ID,
		AuthorChannelID: renderer.AuthorExternalChannelID,
		VideoID:         a.videoID,
		Type:            "super_chat",
		Timestamp:       renderer.TimestampText.SimpleText,
		TimestampUsec:   timestampUsec,
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

	return nil
}

func toHex(d int) string {
	return strconv.FormatInt(int64(d), 16)
}

func retrieveImageURL(thumbnails []chatlog.Thumbnail) string {
	lastIndex := len(thumbnails) - 1
	if lastIndex == -1 {
		return ""
	}
	return thumbnails[lastIndex].URL
}

func parseSuperChat(str string) (string, float64, error) {
	unit := strings.TrimRight(str, "0123456789.,")
	s := strings.TrimLeft(str, unit)
	s = strings.ReplaceAll(s, ",", "")
	amount, err := strconv.ParseFloat(s, 64)
	unit = strings.TrimSpace(strings.ReplaceAll(unit, "￥", "¥"))

	return unit, amount, err
}

func parseMessage(message *chatlog.Message) ([]models.MessageElement, error) {
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
