package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	TELEGRAM_API_URL = "https://api.telegram.org/"
)

func getTelegramBotUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", os.Getenv("TELEGRAM_API_TOKEN"))
}

type SendPhotoObject struct {
	ChatId  string `json:"chat_id"`
	Photo   string `json:"photo"`
	Caption string `json:"caption"`
}

func SendPhoto(imageUrl string, caption string) {
	url := fmt.Sprintf("%s/sendPhoto", getTelegramBotUrl())

	jsonVal, _ := json.Marshal(SendPhotoObject{
		ChatId:  os.Getenv("TELEGRAM_CHAT_ID"),
		Photo:   imageUrl,
		Caption: caption,
	})

	if os.Getenv("LOG_INSTEAD_OF_SEND") == "true" {
		fmt.Printf("jsonVal = %+v\n", jsonVal)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(string(body))
	}
}
