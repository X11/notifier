package main

import (
	"regexp"
	"sync"

	"github.com/mmcdole/gofeed"
)

type FeedNotifier struct {
	feedURL string
}

func (fn *FeedNotifier) getFeed() *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(fn.feedURL)
	if err != nil {
		panic(err)
	}

	return feed
}

type ImageFeedNotifier struct {
	FeedNotifier
	imageRegex *regexp.Regexp
}

func (ifn *ImageFeedNotifier) getImage(target string, i int) string {
	matches := ifn.imageRegex.FindStringSubmatch(target)
	return matches[i]
}

type Execute interface {
	Execute(*sync.WaitGroup)
}

func main() {
	notifiers := []Execute{
		NewXkcd(),
		NewCommitstrip(),
		NewToonhole(),
	}

	var wg sync.WaitGroup

	for _, notifier := range notifiers {
		wg.Add(1)
		go notifier.Execute(&wg)
	}

	wg.Wait()

	if isDirtyState() {
		updateState()
	}
}
