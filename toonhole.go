package main

import (
	"fmt"
	"regexp"

	"github.com/mmcdole/gofeed"
)

type ToonholeState struct {
	Toonhole struct {
		GUID string `json:"guid"`
	} `json:"toonhole"`
}

const (
	TOONHOLE_FEED_URL = "http://toonhole.com/feed/"
)

var (
	TOONHOLE_IMAGE_SOURCE_REGEX = regexp.MustCompile(`img [^>]* src="([^"]*)"`)
)

func getToonholeGUID() string {
	return GetState().Toonhole.GUID
}

func updateToonholeGUID(guid string) {
	GetState().Toonhole.GUID = guid
}

func NotifyToonhole() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(TOONHOLE_FEED_URL)
	if err != nil {
		panic(err)
	}

	item := feed.Items[0]
	if item.GUID != getToonholeGUID() {
		image := TOONHOLE_IMAGE_SOURCE_REGEX.FindStringSubmatch(item.Content)[1]
		fmt.Println("[TOONHOLE] Sending new image: " + image)
		updateToonholeGUID(item.GUID)
		SendPhoto(image, item.Title)
	}
}
