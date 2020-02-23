package gitthing

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// HashTree will "cd" into a GitPath and extract hashes of directories
// currently: "git ls-files -s f ./some/subdir | git hash-object --stdin"
// TODO: don't rely on "sh".. maybe worth benchmarking
func HashTree(GitPath string, dirs []string) {

	for _, subDir := range dirs {
		log.Debugf("%s: Hashing: %s\n", GitPath, subDir)

		// construct the command we'll run.. the "." creates ./somepath
		runme := fmt.Sprintf("git ls-files -s f .%s | git hash-object --stdin", subDir)
		log.Debugf("%s: running: %s\n", GitPath, runme)

		// start the command.. note we'll interpret with "sh"
		cmd := exec.Command("sh", "-c", runme)

		// ensure we run it inside the correct folder (the target git repo)
		// NOTE: there is another possible way to do this using git -C (untested)
		cmd.Dir = GitPath

		data, err := cmd.Output()
		if err != nil {
			log.Debugln(string(data))
			log.Fatal(err)
		}

		// trim newlines etc and get the hash
		hash := strings.TrimSpace(string(data))
		log.Debugf("%s: %s %s\n", GitPath, hash, subDir)

	}

}
