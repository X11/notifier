package main

import (
	"fmt"
	"regexp"

	"github.com/mmcdole/gofeed"
)

type CommitstripState struct {
	Commitstrip struct {
		GUID string `json:"guid"`
	} `json:"commitstrip"`
}

const (
	COMMITSTRIP_FEED_URL = "http://www.commitstrip.com/en/feed/"
)

var (
	COMMITSTRIP_IMAGE_SOURCE_REGEX = regexp.MustCompile(`img src="([^"]*)"`)
)

func getCommitstripGUID() string {
	return GetState().Commitstrip.GUID
}

func updateCommitstripGUID(guid string) {
	GetState().Commitstrip.GUID = guid
}

func NotifyCommitstrip() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(COMMITSTRIP_FEED_URL)
	if err != nil {
		panic(err)
	}

	item := feed.Items[0]
	if item.GUID != getCommitstripGUID() {
		image := COMMITSTRIP_IMAGE_SOURCE_REGEX.FindStringSubmatch(item.Content)[1]
		fmt.Println("[COMMITSTRIP] Sending new image: " + image)
		updateCommitstripGUID(item.GUID)
		SendPhoto(image, item.Title)
	}
}
