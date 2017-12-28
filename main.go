package main

func main() {
	NotifyCommitstrip()
	NotifyXkcd()
	NotifyToonhole()
	if isDirtyState() {
		updateState()
	}
}
