package scrape

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func getStringInBetween(str, start, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

func RetrieveChannelImageURL(channelURL string) (imageURL string, err error) {
	resp, err := http.Get(channelURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	imageURL = getStringInBetween(string(b), `"thumbnailUrl" href="`, `"`)

	return imageURL, nil
}
