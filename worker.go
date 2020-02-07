package gitthing

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

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
	log.Println(gw.repo)
	GitPath := fmt.Sprintf("%s", gw.repo)
	log.Debugf("GitPath: %s", GitPath)
	if _, err := os.Stat(GitPath); os.IsNotExist(err) {
		log.Printf("doing first clone for: %s\n", gw.repo)
		_, err = git.PlainClone(GitPath, false, &git.CloneOptions{
			URL:      gw.repo,
			Auth:     auth,
			Progress: os.Stdout,
		})
		return err
	}

	log.Printf("doing force pull for: %s\n", gw.repo)
	r, err := git.PlainOpen(GitPath)
	if err != nil {
		log.Println(err)
		return
	}

	remotes, err := r.Remotes()
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, remote := range remotes {
		log.Println("remote: ", remote)
		log.Println("Fetching: ", remote.Config().Name, "via", remote.Config().URLs)
		err := remote.Fetch(&git.FetchOptions{
			RemoteName: remote.Config().Name,
			Force:      true,
			Auth:       auth,
			// Progress:   os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			log.Println("Error:", remote.Config().Name, err)
			continue
		}
	}

	branches, err := ListBranches(GitPath)
	for _, b := range branches {
		fmt.Println("detected branch: " + b)
		err := CheckoutBranch(GitPath, b)
		if err != nil {
			log.Fatalln(err)
		}
		ExecTest(GitPath)
		log.Infoln("Sleeping a bit")
		time.Sleep(10 * time.Second)
	}
	return err

}



func ExecTest(GitPath string) {

	cmd := exec.Command("tree")
	cmd.Dir = GitPath
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("---\n%s\n---\n", out)

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
	if err != nil {
		log.Fatal(err)
	}

	refs, err := r.References()
	if err != nil {
		log.Fatal(err)
	}

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

// CheckoutBranch will effectively "cd" into a git repo and do a git checkout
// NOTE.. this will prefix "refs/remotes/origin/" to 'BranchName'
func CheckoutBranch(GitPath string, BranchName string) (err error) {

	//CheckMeOut := fmt.Sprintf("refs/remotes/origin/" + BranchName)
	//r, _ := git.PlainClone(GitPath, false, &git.CloneOptions{
	//	URL: url,
	//})
	r, err := git.PlainOpen(GitPath)
	if err != nil {
		log.Fatal(err)
	}
	w, _ := r.Worktree()

	//err = r.Fetch(&git.FetchOptions{
	//	RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	//	Force: true,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//}
	err = w.Checkout(&git.CheckoutOptions{
		//Branch: plumbing.ReferenceName("refs/heads/" + BranchName),
		//Branch: plumbing.ReferenceName(BranchName),
		Branch: plumbing.ReferenceName("refs/remotes/origin/" + BranchName),
		Force: true,
		Keep: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("[%s] checked out branch: %s\n", GitPath, BranchName)
	return err
}
