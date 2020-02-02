package gitthing

import "time"

type Config struct {
	Global Global `yaml:"global",yaml:"global"`
}

type Global struct {
	PollFrequency time.Duration `yaml:"pollFrequency"`
	PublicAddress string        `yaml:"publicAddress"`
	SshKeyPath    string        `yaml:"sshKeyPath"`
}

type Repo struct {
	Url        string `yaml:"url"`
	Path       string `yaml:"path"`
	Ref        string `yaml:"ref"`
	SshKeyPath string `yaml:"sshKeyPath"`
}

type Provider struct {
	Name  string `yaml:"name"`
	Match string `yaml:"match"`
	//PollFreq time.Duration `yaml:"pollFreq"` // TODO - implement later
}
