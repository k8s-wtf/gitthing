package gitthing

import (
	"time"

	log "github.com/sirupsen/logrus"
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
			r.Trigger()
		case c := <-r.stop:
			c <- nil
		}
	}
}

func (r *RecLoop) Trigger() {
	log.Warnln("tick")
	for _, job := range r.repos {
		log.Debugln(job)
		if job.PollFrequency == 0*time.Second {
			// omitempty will default to time.Duration(0 * time.Second)
			// Assume Global default
			job.PollFrequency = r.pollFreq
		}

		log.Debugf("NewGitWorker: (%v)\n", job)
		w := NewGitWorker(
			job.SshKeyPath,
			job.Url,
			"",
			job.PollFrequency,
		)
		err := w.Do()
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}

func (r *RecLoop) Stop() error {
	stopCh := make(chan error)
	r.stop <- stopCh
	return <-stopCh
}
