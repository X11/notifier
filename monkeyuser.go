package main

import (
	"fmt"
	"regexp"
	"sync"
)

type MonkeyuserState struct {
	Monkeyuser struct {
		GUID string `json:"guid"`
	} `json:"toonhole"`
}

func getMonkeyuserGUID() string {
	return GetState().Monkeyuser.GUID
}

func updateMonkeyuserGUID(guid string) {
	GetState().Monkeyuser.GUID = guid
	SetStateDirty()
}

type Monkeyuser struct {
	ImageFeedNotifier
}

func NewMonkeyuser() *Monkeyuser {
	n := &Monkeyuser{}
	n.feedURL = "http://www.monkeyuser.com/feed.xml"
	n.imageRegex = regexp.MustCompile(`img[^>]* src="([^"]*)"`)
	return n
}

func (x *Monkeyuser) Execute(wg *sync.WaitGroup) {
	defer wg.Done()

	feed := x.getFeed()
	item := feed.Items[0]
	if item.GUID != getMonkeyuserGUID() {
		image := x.getImage(item.Description, 1)
		fmt.Println("[TOONHOLE] Sending new image: " + image)
		updateMonkeyuserGUID(item.GUID)
		SendPhoto(image, item.Title)
	}
}
