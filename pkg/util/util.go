package util

import (
	"bytes"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// ExpandPath expands tilda(~) to absoulte path
func ExpandPath(path string) string {
	usr, _ := user.Current()
	home := usr.HomeDir

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	if strings.HasPrefix(path, "./") {
		path = filepath.Join(pwd, path[2:])
	}

	return path
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ExecutableExists reports whether the named executable exists in $PATH
func ExecutableExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// RunCommand runs given command on host machine
func RunCommand(cmd *exec.Cmd) (string, error) {
	// set var to get the output
	var out bytes.Buffer

	// set the output to our variable
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
