package main

import (
	"fmt"
	"github.com/k8s-wtf/gitthing"
	"time"
)

func main() {
	r := gitthing.NewRecLoop(time.Second * 2, make([]gitthing.Repo,0))
	go r.Run()
	time.Sleep(time.Second * 3)
	fmt.Printf("%s", r.Stop())
}
