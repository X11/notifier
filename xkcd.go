package main

import (
	"fmt"
	"regexp"

	"github.com/mmcdole/gofeed"
)

type XkcdState struct {
	Xkcd struct {
		GUID string `json:"guid"`
	} `json:"xkcd"`
}

const (
	XKCD_FEED_URL = "https://xkcd.com/rss.xml"
)

var (
	XKCD_SOURCE_REGEX = regexp.MustCompile(`img src="([^"]*)" title="([^"]*)"`)
)

func getXkcdGUID() string {
	return GetState().Xkcd.GUID
}

func updateXkcdGUID(guid string) {
	GetState().Xkcd.GUID = guid
}

func NotifyXkcd() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(XKCD_FEED_URL)
	if err != nil {
		panic(err)
	}

	item := feed.Items[0]
	if item.GUID != getXkcdGUID() {
		matches := XKCD_SOURCE_REGEX.FindStringSubmatch(item.Description)
		image := matches[1]
		caption := fmt.Sprintf("%s\n%s", item.Title, matches[2])
		fmt.Println("[XKCD] Sending new image: " + image)
		updateXkcdGUID(item.GUID)
		SendPhoto(image, caption)
	}
}
