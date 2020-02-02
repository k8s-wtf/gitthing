package gitthing

import (
	"time"
)

type RecLoop struct {
	pollFreq time.Duration
	repos    []Repo
	stop     chan chan error
}

func NewRecLoop(pollFreq time.Duration, repos []Repo) *RecLoop {
	return &RecLoop{pollFreq: pollFreq, repos: repos, stop: make(chan chan error)}
}

func (r *RecLoop) Run() {
	ticker := time.NewTicker(r.pollFreq)
	for {
		select {
		case <-ticker.C:
			println("tick")
		case c := <-r.stop:
			c <- nil
		}
	}
}

func (r *RecLoop) Stop() error {
	stopCh := make(chan error)
	r.stop <- stopCh
	return <-stopCh
}
