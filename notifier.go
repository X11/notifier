package main

import (
	"regexp"

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
	if len(matches) >= i {
		return matches[i]
	}
	return ""
}
