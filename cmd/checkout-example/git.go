package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)


func main(){
	sshKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa-gitthing"
	pem, _ := ioutil.ReadFile(sshKeyPath)
	signer, _ := ssh.ParsePrivateKey(pem)
	auth := &ssh2.PublicKeys{User: "git", Signer: signer}


	url := "git@github.com:k8s-wtf/gitthing-example.git"
	_, err := git.PlainClone(url, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     auth,
	})
	if err != nil {
		fmt.Printf("pfff")
	}
}
