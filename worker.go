package gitthing

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"os"
)

type GitWorker struct {
	sshKeyPath    string
	repo          string
	branchPattern string
}

func NewGitWorker(sshKeyPath string, repo string, branchPattern string) *GitWorker {
	return &GitWorker{sshKeyPath: sshKeyPath, repo: repo, branchPattern: branchPattern}
}

func (gw *GitWorker) Do() (err error) {
	pem, err := ioutil.ReadFile(gw.sshKeyPath)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(pem)
	if err != nil {
		return err
	}
	auth := &ssh2.PublicKeys{
		User:                  "git",
		Signer:                signer,
	}
	_, err = git.PlainClone(gw.repo, false, &git.CloneOptions{
		URL:      gw.repo,
		Auth:     auth,
		Progress: os.Stdout,
	})
	return err
}