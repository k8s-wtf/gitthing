package main

import (
	"io/ioutil"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/k8s-wtf/gitthing"
	log "github.com/sirupsen/logrus"
)

const configPath = "example-config.yaml"

func main() {
	log.SetLevel(log.DebugLevel)
	configRaw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("opening configfile: %s", err)
	}
	config := &gitthing.Config{}
	err = yaml.Unmarshal(configRaw, config)
	if err != nil {
		log.Fatalf("parsing config: %s", err)
	}

	for _, r := range config.Repos {

		if r.PollFrequency == 0 * time.Second {
			// omitempty will default to time.Duration(0 * time.Second)
			// Assume Global default
			r.PollFrequency = config.Global.PollFrequency
		}

		log.Debugf("NewGitWorker: (%v)\n", r)
		w := gitthing.NewGitWorker(
			r.SshKeyPath,
			r.Url,
			"",
			r.PollFrequency,
		)
		err := w.Do()
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	r := gitthing.NewRecLoop(config.Global.PollFrequency, make([]gitthing.Repo, 0))
	go r.Run()

	time.Sleep(time.Second * 300)
	// err = r.Stop()
	// if err != nil {
	// 	log.Errorf("shutting down loop: %s", err)
	// }
	// fmt.Printf("%s", r.Stop())

}
