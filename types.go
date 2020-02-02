package gitthing

import "time"

type Config struct {
	PollFreq time.Duration
}

type Repos []Repo

type Repo struct {
	Url        string `json:"url"`
	Path       string `json:"path"`
	Ref        string `json:"ref"`
	SshKeyPath string `json:"sshKeyPath"`
}

type Provider struct {
	Name  string `json:"name"`
	Match string `json:"match"`
	PollFreq time.Duration `json:"pollFreq"`
}
