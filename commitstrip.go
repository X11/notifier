package main

import (
	"fmt"
	"regexp"
	"sync"
)

type CommitstripState struct {
	Commitstrip struct {
		GUID string `json:"guid"`
	} `json:"commitstrip"`
}

func getCommitstripGUID() string {
	return GetState().Commitstrip.GUID
}

func updateCommitstripGUID(guid string) {
	GetState().Commitstrip.GUID = guid
}

type Commitstrip struct {
	ImageFeedNotifier
}

func NewCommitstrip() *Commitstrip {
	n := &Commitstrip{}
	n.feedURL = "http://www.commitstrip.com/en/feed/"
	n.imageRegex = regexp.MustCompile(`img src="([^"]*)"`)
	return n
}

func (x *Commitstrip) Execute(wg *sync.WaitGroup) {
	defer wg.Done()

	feed := x.getFeed()
	item := feed.Items[0]
	if item.GUID != getCommitstripGUID() {
		image := x.getImage(item.Content, 1)
		fmt.Println("[COMMITSTRIP] Sending new image: " + image)
		updateCommitstripGUID(item.GUID)
		SendPhoto(image, item.Title)
	}
}
