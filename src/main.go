package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func printUsage() {
	fmt.Printf("Usage: %s <path-to-python-project>\n", os.Args[0])
}

func expandHome(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	return path
}

// Exists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// as util
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func run(cmd *exec.Cmd) (string, error) {

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

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	base := filepath.Base(source)

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if source == path {
				return nil
			}
			path += "/"
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = path[len(base)+1:]
		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		return err
	}
	if err = archive.Flush(); err != nil {
		return err
	}
	return nil
}

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	projPath := expandHome(os.Args[1])
	target := fmt.Sprintf("%s.%s", projPath, "zip")
	logrus.Infof("Python Project: %s\n", projPath)
	logrus.Infof("Output: %s\n", target)

	if fileExists(target) {
		logrus.Fatal("file already exists at ", target)
	}

	var pip string
	if commandExists("pip3") {
		pip = "pip3"
	} else if commandExists("pip") {
		pip = "pip"
	} else {
		logrus.Fatal("'pip3' or 'pip' must be installed")
	}

	// install dependencies
	logrus.Info("checking requirements.txt in python project")
	requirementsPath := filepath.Join(projPath, "requirements.txt")
	if fileExists(requirementsPath) {
		logrus.Info("requirements.txt exists")
	} else {
		logrus.Infof("%s not found. create a new one.\n", requirementsPath)

		pipreqs := "pipreqs"
		if !commandExists(pipreqs) {
			logrus.Error("'pipreqs' not installed")
			logrus.Fatalf("please run '%s install pipreqs' to install dependency", pip)
		}

		cmd := exec.Command(pipreqs, projPath, "--savepath", requirementsPath)
		out, err := run(cmd)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info(out)
	}

	logrus.Info("install python libraries")
	cmd := exec.Command(pip, "install", "-r", requirementsPath, "-t", projPath)
	out, err := run(cmd)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(out)

	logrus.Infof("zip %s to %s", projPath, target)
	err = zipit(projPath, target)
	if err != nil {
		logrus.Fatal("Failed to zip python project. ", err)
	}
}
