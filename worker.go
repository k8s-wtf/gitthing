package gitthing

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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
		User:   "git",
		Signer: signer,
	}
	fmt.Println(gw.repo)
	GitPath := fmt.Sprintf("%s\n", gw.repo)
	if _, err := os.Stat(GitPath); os.IsNotExist(err) {
		fmt.Printf("doing first clone for: %s\n", gw.repo)
		_, err = git.PlainClone(GitPath, false, &git.CloneOptions{
			URL:      gw.repo,
			Auth:     auth,
			Progress: os.Stdout,
		})
		return err
	}

	fmt.Printf("doing force pull for: %s\n", gw.repo)
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
		fmt.Println("remote: ", remote)
		fmt.Println("Fetching: ", remote.Config().Name, "via", remote.Config().URLs)
		err := remote.Fetch(&git.FetchOptions{
			RemoteName: remote.Config().Name,
			Force:      true,
			Auth:       auth,
			// Progress:   os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			fmt.Println("Error:", remote.Config().Name, err)
			continue
		}
	}

	return err

}

// ListBranches returns a list of branchs for a locally checkedout git (and an err)
// accepts just one variables.. the GitPath where the repo is checked out
// it will filter to only show repo's matching the prefix: "refs/remotes/origin/"
// EG usage
//
//	branches, err := ListBranches(GitPath)
//	for _, b := range branches {
//		fmt.Println("detected branch: " + b)
//	}

func ListBranches(GitPath string) (branchList []string, err error) {

	r, err := git.PlainOpen(GitPath)
	CheckIfError(err)

	refs, err := r.References()
	CheckIfError(err)

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		branchName := ref.Name().String()
		if strings.HasPrefix(branchName, "refs/remotes/origin/") {
			// Do we want human readable variable here or is it gonna be a ball ache?
			branchNameHuman := strings.Replace(branchName, "refs/remotes/origin/", "", -1)
			branchList = append(branchList, branchNameHuman)
		}
		return nil
	})
	return branchList, err

}
