package gitthing

import (
	"fmt"
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
	fmt.Println(gw.repo)
	GitPath := fmt.Sprintf("%s\n", gw.repo)
	if _, err := os.Stat(GitPath); os.IsNotExist(err) {
		fmt.Sprintln("doing first clone for: %s", gw.repo)
		_, err = git.PlainClone(GitPath, false, &git.CloneOptions{
			URL:      gw.repo,
			Auth:     auth,
			Progress: os.Stdout,
		})
		return err
	}


	fmt.Sprintln("doing force pull for: %s", gw.repo)
	r, err := git.PlainOpen(GitPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	remotes, err := r.Remotes()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, remote := range remotes {
		fmt.Println("Fetching: ", remote.Config().Name, "via", remote.Config().URLs)
		err := remote.Fetch(&git.FetchOptions{
			RemoteName: remote.Config().Name,
			Force:      true,
			Auth:     auth,
			Progress: os.Stdout,

		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			fmt.Println("Error:", remote.Config().Name, err)
			continue
		}
	}
	return err

}