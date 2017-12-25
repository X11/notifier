package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	GITHUB_API_URL = "https://api.github.com"
)

var (
	FETCHED_STATE *State
)

type Gist struct {
	Files map[string]GistFile `json:"files"`
}

type GistFile struct {
	Content string `json:"content"`
}

func GetState() *State {
	if FETCHED_STATE != nil {
		return FETCHED_STATE
	}

	return fetchState()
}

func getGistUrl() string {
	return fmt.Sprintf("%s/gists/%s", GITHUB_API_URL, os.Getenv("GITHUB_GIST_ID"))
}

func getGistFileContent(g Gist) string {
	return g.Files[os.Getenv("GITHUB_GIST_FILE")].Content
}

func getStateFromGist(g Gist) *State {
	FETCHED_STATE = &State{}
	if err := json.Unmarshal([]byte(getGistFileContent(g)), &FETCHED_STATE); err != nil {
		panic(err)
	}
	return FETCHED_STATE
}

func fetchState() *State {
	resp, err := http.Get(getGistUrl())
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(string(body))
	}

	b, err := ioutil.ReadAll(resp.Body)
	gist := Gist{}

	if err = json.Unmarshal(b, &gist); err != nil {
		panic(err)
	}

	return getStateFromGist(gist)
}

func marshalState() []byte {
	val, err := json.Marshal(FETCHED_STATE)
	if err != nil {
		panic(err)
	}
	return val
}

func prepareUpdatedGist() []byte {
	val, err := json.Marshal(Gist{
		Files: map[string]GistFile{
			os.Getenv("GITHUB_GIST_FILE"): {Content: string(marshalState())},
		},
	})
	if err != nil {
		panic(err)
	}
	return val
}

func updateState() {
	req, err := http.NewRequest("PATCH", getGistUrl(), bytes.NewBuffer(prepareUpdatedGist()))
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
