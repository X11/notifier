package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	GITHUB_API_URL = "https://api.github.com"
)

var (
	RAW_STATE     string
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
	RAW_STATE = getGistFileContent(g)
	FETCHED_STATE = &State{}
	if err := json.Unmarshal([]byte(RAW_STATE), &FETCHED_STATE); err != nil {
		panic(err)
	}
	return FETCHED_STATE
}

func newGithubRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(os.Getenv("GITHUB_AUTHENTICATION"))))
	return req, err
}

func fetchState() *State {
	req, err := newGithubRequest("GET", getGistUrl(), nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
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

func isDirtyState() bool {
	return RAW_STATE != string(marshalState())
}

func updateState() {
	req, err := newGithubRequest("PATCH", getGistUrl(), bytes.NewBuffer(prepareUpdatedGist()))
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
