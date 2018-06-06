package main

import (
	"os"

	gstate "github.com/X11/go-gist-store"
)

type State struct {
	CommitstripState
	XkcdState
	ToonholeState
	MonkeyuserState
}

var s *State
var gs *gstate.GState
var dirty = false

func FetchState() {
	gs = gstate.New(os.Getenv("GITHUB_GIST_ID"), os.Getenv("GITHUB_GIST_FILE"), os.Getenv("GITHUB_AUTHENTICATION"))
	s = &State{}
	gs.Get(s)
}

func SetStateDirty() {
	dirty = true
}

func GetState() *State {
	return s
}

func UpdateState() {
	if dirty {
		gs.Update(s)
	}
}
