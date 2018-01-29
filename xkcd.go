package main

import (
	"fmt"
	"regexp"
	"sync"
)

type XkcdState struct {
	Xkcd struct {
		GUID string `json:"guid"`
	} `json:"xkcd"`
}

func getXkcdGUID() string {
	return GetState().Xkcd.GUID
}

func updateXkcdGUID(guid string) {
	GetState().Xkcd.GUID = guid
}

type Xkcd struct {
	ImageFeedNotifier
}

func NewXkcd() *Xkcd {
	n := &Xkcd{}
	n.feedURL = "https://xkcd.com/rss.xml"
	n.imageRegex = regexp.MustCompile(`img src="([^"]*)" title="([^"]*)"`)
	return n
}

func (x *Xkcd) Execute(wg *sync.WaitGroup) {
	defer wg.Done()

	feed := x.getFeed()
	item := feed.Items[0]
	if item.GUID != getXkcdGUID() {
		image := x.getImage(item.Description, 1)
		alt := x.getImage(item.Description, 2)
		caption := fmt.Sprintf("%s\n%s", item.Title, alt)
		fmt.Println("[XKCD] Sending new image: " + image)
		updateXkcdGUID(item.GUID)
		SendPhoto(image, caption)
	}
}
