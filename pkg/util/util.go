package util

import (
	"bufio"
	"bytes"
	"fmt"
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
func RunCommand(name string, args ...string) (string, error) {
	command := exec.Command(name, args...)

	logrus.Info("execute command: " + command.String())

	// set var to get the output
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// set the output to our variable
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		return "", fmt.Errorf("%s. %s", stderr, err)
	}

	return stdout.String(), nil
}

func Write(path string, contents ...string) error {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("failed to open file %s. %s", path, err)
	}

	datawriter := bufio.NewWriter(file)

	for _, c := range contents {
		_, err = datawriter.WriteString(c + "\n")

		if err != nil {
			return fmt.Errorf("failed to write to file %s. %s", path, err)
		}
	}

	datawriter.Flush()
	file.Close()

	return nil
}
