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


	r := gitthing.NewRecLoop(config.Global.PollFrequency, config.Repos)
	go r.Run()

	time.Sleep(time.Second * 300)
	err = r.Stop()
	if err != nil {
		log.Errorf("shutting down loop: %s", err)
	}
	log.Infof("%s", r.Stop())

}
