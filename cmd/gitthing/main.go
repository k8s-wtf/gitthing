package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/k8s-wtf/gitthing"
	"io/ioutil"
	"log"
	"time"
)

const configPath = "example-config.yaml"

func main() {
	configRaw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("opening configfile: %s", err)
	}
	config := &gitthing.Config{}
	err = yaml.Unmarshal(configRaw, config)
	if err != nil {
		log.Fatalf("parsing config: %s", err)
	}
	fmt.Printf("%+v", config)

	r := gitthing.NewRecLoop(config.Global.PollFrequency, make([]gitthing.Repo,0))
	go r.Run()
	time.Sleep(time.Second * 3)
	fmt.Printf("%s", r.Stop())
}
