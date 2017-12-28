package main

func main() {
	NotifyCommitstrip()
	NotifyXkcd()
	if isDirtyState() {
		updateState()
	}
}
