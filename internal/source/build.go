package source

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Build builds the source function to .so file.
func Build(appPath string, clean bool) (string, error) {
	// check if the file exists
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return "", fmt.Errorf("the file %s doesn't exist", appPath)
	}

	// build
	version := runtime.GOOS
	dir, _ := filepath.Split(appPath)
	so := dir + "source.so"

	// clean build
	if clean {
		// .so file exists, remove it.
		if _, err := os.Stat(so); !os.IsNotExist(err) {
			err = os.Remove(so)
			if err != nil {
				return "", fmt.Errorf("clean build the file %s failed", appPath)
			}
		}
	}

	if version == "linux" {
		cmd := exec.Command("/bin/sh", "-c", "CGO_ENABLED=1 GOOS=linux go build -buildmode=plugin -o "+so+" "+appPath)
		err := cmd.Start()
		if err != nil {
			return "", err
		}
		err = cmd.Wait()
		return so, err
	} else if version == "darwin" {
		cmd := exec.Command("/bin/sh", "-c", "go build -buildmode=plugin -o "+so+" "+appPath)
		err := cmd.Start()
		if err != nil {
			return "", err
		}
		err = cmd.Wait()
		return so, err
	} else {
		return "", errors.New("Not Implemented")
	}

}
