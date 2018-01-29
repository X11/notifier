package main

import (
	"fmt"
	"regexp"
	"sync"
)

type ToonholeState struct {
	Toonhole struct {
		GUID string `json:"guid"`
	} `json:"toonhole"`
}

func getToonholeGUID() string {
	return GetState().Toonhole.GUID
}

func updateToonholeGUID(guid string) {
	GetState().Toonhole.GUID = guid
	SetStateDirty()
}

type Toonhole struct {
	ImageFeedNotifier
}

func NewToonhole() *Toonhole {
	n := &Toonhole{}
	n.feedURL = "http://toonhole.com/feed/"
	n.imageRegex = regexp.MustCompile(`img [^>]* src="([^"]*)"`)
	return n
}

func (x *Toonhole) Execute(wg *sync.WaitGroup) {
	defer wg.Done()

	feed := x.getFeed()
	item := feed.Items[0]
	if item.GUID != getToonholeGUID() {
		image := x.getImage(item.Content, 1)
		fmt.Println("[TOONHOLE] Sending new image: " + image)
		updateToonholeGUID(item.GUID)
		SendPhoto(image, item.Title)
	}
}
