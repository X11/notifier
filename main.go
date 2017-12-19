package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/mmcdole/gofeed"
)

func sendPhoto(imageUrl string, caption string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", os.Getenv("TELEGRAM_API_TOKEN"))

	body := struct {
		ChatId  string `json:"chat_id"`
		Photo   string `json:"photo"`
		Caption string `json:"caption"`
	}{
		ChatId:  os.Getenv("TELEGRAM_CHAT_ID"),
		Photo:   imageUrl,
		Caption: caption,
	}

	jsonVal, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic("Response is not 200")
	}
}

func getCommitStripGuid() string {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/gists/%s", os.Getenv("GITHUB_GIST_ID")))
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)

	body := struct {
		Files map[string]struct {
			Content string
		}
	}{}

	err = json.Unmarshal(b, &body)
	if err != nil {
		panic(err)
	}

	return body.Files["notifier-commitstrip-guid"].Content
}

func patchCommitStripGuid(guid string) {
	url := fmt.Sprintf("https://api.github.com/gists/%s", os.Getenv("GITHUB_GIST_ID"))

	jsonVal, _ := json.Marshal(struct {
		Files map[string]struct {
			Content string `json:"content"`
		} `json:"files"`
	}{
		Files: map[string]struct {
			Content string `json:"content"`
		}{
			"notifier-commitstrip-guid": {Content: guid},
		},
	})

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonVal))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(os.Getenv("GITHUB_AUTHENTICATION"))))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(string(body))
	}

}

func main() {
	guid := getCommitStripGuid()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("http://www.commitstrip.com/en/feed/")
	if err != nil {
		panic(err)
	}

	imageSrc := regexp.MustCompile(`src="([^"]*)"`)

	item := feed.Items[0]
	if item.GUID != guid {
		image := imageSrc.FindStringSubmatch(item.Content)[1]
		sendPhoto(image, item.Title)
		patchCommitStripGuid(item.GUID)
		fmt.Println("Sending new image: " + image)
	}
}
