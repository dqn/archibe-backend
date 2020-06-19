package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dqn/tubekids/dbexec"
	"github.com/dqn/tubekids/models"
	"github.com/dqn/tubekids/youtube/archive"
	"github.com/dqn/tubekids/youtube/scrape"
	"github.com/dqn/ytvi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func channelExists(channels []models.Channel, channelID string) bool {
	for i := range channels {
		if channels[i].ChannelID == channelID {
			return true
		}
	}
	return false
}

func run() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("invalid arguments")
	}

	dsn := os.Args[1]
	videoID := os.Args[2]

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	dbx := dbexec.NewExecutor(db.MustBegin())

	fmt.Println("start fetching chats...")

	fetcher := archive.NewFetcher(videoID)
	acv, err := fetcher.Fetch()
	if err != nil {
		return err
	}

	pr, err := ytvi.GetVideoInfo(videoID)
	if err != nil {
		return err
	}

	pmr := pr.Microformat.PlayerMicroformatRenderer
	if !channelExists(acv.Channels, pmr.ExternalChannelID) {
		channelImageURL, err := scrape.RetrieveChannelImageURL(pmr.OwnerProfileURL)
		if err != nil {
			return err
		}

		ownerChannel := models.Channel{
			ChannelID: pmr.ExternalChannelID,
			Name:      pmr.OwnerChannelName,
			ImageURL:  channelImageURL,
		}

		acv.Channels = append(acv.Channels, ownerChannel)
	}

	fmt.Println("start inserting to database...")

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
	_, err = dbx.Channels.InsertMany(acv.Channels)
	if err != nil {
		return err
	}
	_, err = dbx.Chats.InsertMany(acv.Chats)
	if err != nil {
		return err
	}
	_, err = dbx.Badges.InsertMany(acv.Badges)
	if err != nil {
		return err
	}

	if err = dbx.Tx.Commit(); err != nil {
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
