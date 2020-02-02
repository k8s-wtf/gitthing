package gitthing

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"os"
)

type GitWorker struct {
	cloneOptions  *git.CloneOptions
	checkoutPath  string
	repo          string
	branchPattern string
}

func NewGitWorker(sshKeyPath string, repo string, branchPattern string) (*GitWorker, error) {
	pem, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(pem)
	if err != nil {
		return nil, err
	}
	auth := &ssh2.PublicKeys{
		User:   "git",
		Signer: signer,
	}
	cloneOpts := &git.CloneOptions{
		URL:      repo,
		Auth:     auth,
		Progress: os.Stdout,
	}
	return &GitWorker{cloneOptions: cloneOpts, repo: repo, branchPattern: branchPattern}, nil
}

func (gw *GitWorker) Clone() (err error) {
	_, err = git.PlainClone("foo", false, gw.cloneOptions)
	return err
}

func (gw *GitWorker) Fetch() error {
	r, err := git.PlainOpen("foo")
	if err != nil {
		return err
	}
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "",
		RefSpecs:   nil,
		Depth:      0,
		Auth:       nil,
		Progress:   nil,
		Tags:       0,
		Force:      false,
	})
	return err
}
