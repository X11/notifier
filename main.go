package main

import "sync"

type Execute interface {
	Execute(*sync.WaitGroup)
}

func main() {
	FetchState()

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

	UpdateState()
}
