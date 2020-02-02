package main

import (
	"github.com/go-yaml/yaml"
	"github.com/k8s-wtf/gitthing"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

const configPath = "example-config.yaml"

func main() {
	configRaw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("opening config file: %s", err)
	}
	config := &gitthing.Config{}
	err = yaml.Unmarshal(configRaw, config)
	if err != nil {
		log.Fatalf("parsing config: %s", err)
	}

	for _, i := range config.Repos {
		w, err := gitthing.NewGitWorker(i.SshKeyPath, i.Url, "")
		if err != nil {
			log.Fatal(err)
		}
		err = w.Clone()
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	//r := gitthing.NewRecLoop(config.Global.PollFrequency, make([]gitthing.Repo,0))
	//go r.Run()
	////time.Sleep(time.Second * 3)
	//err = r.Stop()
	//if err != nil {
	//	log.Errorf("shutting down loop: %s", err)
	//}
	//fmt.Printf("%s", r.Stop())
}
