package gitthing

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Stolen from https://gist.github.com/mustafaydemir/c90db8fcefeb4eb89696e6ccb5b28685
// modified to only return subdirectories and files (not including the name of the parent)
func scanRecursive(dirPath string, ignore []string) (folders []string, files []string) {
	// Scan
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		_continue := false
		// Loop : Ignore Files & Folders
		for _, i := range ignore {
			// If ignored path
			if strings.Index(path, i) != -1 {
				// Continue
				_continue = true
			}
		}
		if _continue == false {
			f, err = os.Stat(path)
			// If no error
			if err != nil {
				log.Fatal(err)
			}
			// File & Folder Mode
			fMode := f.Mode()

			// Filter out the top level directory.. we'll already know that name when calling this func
			path := strings.Replace(path, dirPath, "", 1)
			if path != "" {
				// strip leading /
				// path := strings.Replace(path, "/", "", 1)
				// Is folder
				if fMode.IsDir() {
					folders = append(folders, path)
					// Is file
				} else if fMode.IsRegular() {
					// Append to Files Array
					files = append(files, path)
				}
			}
		}
		return nil
	})
	return folders, files
}

// PathIndex will index a path and return a list of subdirectories
func PathIndex(path string) (result []string) {

	IgnoreMe := []string{
		"/.git",
	}

	folders, _ := scanRecursive(path, IgnoreMe)
	for _, i := range folders {
		// log.Debugf("found dir: %s\n", i)
		result = append(result, i)
	}
	return result
}
