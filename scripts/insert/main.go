package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/dqn/tubekids/dbexec"
	"github.com/dqn/tubekids/models"
	"github.com/dqn/tubekids/youtube/archive"
	"github.com/dqn/tubekids/youtube/scrape"
	"github.com/dqn/ytcv"
	"github.com/dqn/ytvi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ArchiveData struct {
	Video    *models.Video    `json:"video"`
	Channels []models.Channel `json:"channels"`
	Chats    []models.Chat    `json:"chats"`
	Badges   []models.Badge   `json:"badges"`
}

func videoExists(db *sqlx.DB, videoID string) (bool, error) {
	videoIDs := []string{}
	if err := db.Select(&videoIDs, "SELECT video_id FROM videos WHERE video_id = $1", videoID); err != nil {
		return false, err
	}
	return (len(videoIDs) == 1), nil
}

func channelExists(channels []models.Channel, channelID string) bool {
	for i := range channels {
		if channels[i].ChannelID == channelID {
			return true
		}
	}
	return false
}

func fetchArchiveData(videoID string) (*ArchiveData, error) {
	pr, err := ytvi.GetVideoInfo(videoID)
	if err != nil {
		return nil, err
	}

	if !pr.VideoDetails.IsLiveContent {
		err = fmt.Errorf("not a live content")
		return nil, err
	}

	fetcher := archive.NewFetcher(videoID)
	acv, err := fetcher.Fetch()
	if err != nil {
		return nil, err
	}

	pmr := pr.Microformat.PlayerMicroformatRenderer
	if !channelExists(acv.Channels, pmr.ExternalChannelID) {
		channelImageURL, err := scrape.RetrieveChannelImageURL(pmr.OwnerProfileURL)
		if err != nil {
			return nil, err
		}

		ownerChannel := models.Channel{
			ChannelID: pmr.ExternalChannelID,
			Name:      pmr.OwnerChannelName,
			ImageURL:  channelImageURL,
		}

		acv.Channels = append(acv.Channels, ownerChannel)
	}

	lengthSeconds, err := strconv.ParseInt(pmr.LengthSeconds, 10, 64)
	if err != nil {
		return nil, err
	}
	viewCount, err := strconv.ParseInt(pmr.ViewCount, 10, 64)
	if err != nil {
		return nil, err
	}
	publishDate, err := time.Parse("2006-01-02", pmr.PublishDate)
	if err != nil {
		return nil, err
	}
	uploadDate, err := time.Parse("2006-01-02", pmr.UploadDate)
	if err != nil {
		return nil, err
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

	a := ArchiveData{
		Video:    &video,
		Channels: acv.Channels,
		Chats:    acv.Chats,
		Badges:   acv.Badges,
	}

	return &a, nil
}

func saveArchiveData(path string, data *ArchiveData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func insertByJSON(db *sqlx.DB, dirname string) error {
	entries, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	dbx := dbexec.NewExecutor(db)
	for i, entry := range entries {
		f, err := os.Open(filepath.Join(dirname, entry.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", entry.Name(), err)
			continue
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", entry.Name(), err)
			continue
		}

		var data ArchiveData
		if err = json.Unmarshal(b, &data); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", entry.Name(), err)
			continue
		}

		fmt.Printf("%s: start inserting...\n", data.Video.VideoID)

		if err := insertArchiveData(*dbx, &data); err != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to insert: %s\n", data.Video.VideoID, err)
			continue
		}

		remaining := len(entries) - i
		fmt.Printf("%s: finished inserting (remaining: %d) \n", data.Video.VideoID, remaining)
	}

	return nil
}

func insertArchiveData(dbx dbexec.DBExecutor, data *ArchiveData) error {
	_, err := dbx.Videos.InsertOne(data.Video)
	if err != nil {
		return err
	}
	_, err = dbx.Channels.InsertMany(data.Channels)
	if err != nil {
		return err
	}
	_, err = dbx.Chats.InsertMany(data.Chats)
	if err != nil {
		return err
	}
	_, err = dbx.Badges.InsertMany(data.Badges)
	if err != nil {
		return err
	}

	return nil
}

func run() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("invalid arguments")
	}

	dsn := os.Args[1]
	channelID := os.Args[2]

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// insertByJSON(db, "./data/")
	// return nil

	fmt.Printf("%s: start fetching channel videos...\n", channelID)

	videos, err := ytcv.FetchAll(channelID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	ch := make(chan *ArchiveData, len(videos))

	wg.Add(len(videos))
	for _, video := range videos {
		go func(videoID string) {
			defer wg.Done()

			fmt.Printf("%s: start fetching chats...\n", videoID)

			exists, err := videoExists(db, videoID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", videoID, err)
				return
			}
			if exists {
				fmt.Printf("%s: already fetched\n", videoID)
				return
			}

			a, err := fetchArchiveData(videoID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", videoID, err)
				return
			}

			fmt.Printf("%s: finished fetching\n", a.Video.VideoID)
			ch <- a
		}(video.VideoID)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	dbx := dbexec.NewExecutor(db)

	for data := range ch {
		fmt.Printf("%s: start inserting...\n", data.Video.VideoID)

		if err = insertArchiveData(*dbx, data); err != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to insert: %s\n", data.Video.VideoID, err)

			path := fmt.Sprintf("./data/%s.json", data.Video.VideoID)
			if err = saveArchiveData(path, data); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", data.Video.VideoID, err)
			} else {
				fmt.Printf("saved: %s\n", path)
			}
			continue
		}

		fmt.Printf("%s: finished inserting \n", data.Video.VideoID)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
