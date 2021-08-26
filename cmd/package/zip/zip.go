/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BlueWhaleKo/python-lambda-zipper/pkg/archive"
	"github.com/BlueWhaleKo/python-lambda-zipper/pkg/util"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type zipArgs struct {
	ProjectPath string
	OutputPath  string
}

var zipargs = &zipArgs{}

var zipCmd = &cobra.Command{
	Use:   "zip --[flags] [options]",
	Short: "pack python projects into a zip file",
	Long:  "pack python projects into a zip file",
	Run: func(cmd *cobra.Command, args []string) {
		main(cmd, args)
	},
}

func NewZipCommand() *cobra.Command {
	return zipCmd
}

func init() {
	// parse args
	zipCmd.Flags().StringVar(&zipargs.ProjectPath, "project-path", "", "(required) path to python project")
	zipCmd.Flags().StringVar(&zipargs.OutputPath, "output-path", "", "(required) path to output zip file")

	zipCmd.MarkFlagRequired("project-path")
	zipCmd.MarkFlagRequired("output-path")
}

func main(cmd *cobra.Command, args []string) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	projDir := util.ExpandPath(zipargs.ProjectPath)
	target := util.ExpandPath(zipargs.OutputPath)
	logrus.Infof("Python Project: %s\n", projDir)
	logrus.Infof("Output: %s\n", target)
	logrus.Infof("Tmp Dir: %s\n", tmpDir)

	if err := copy.Copy(projDir, tmpDir); err != nil {
		log.Fatal(err)
	}

	pip := ""
	if util.ExecutableExists("pip3") {
		pip = "pip3"
	} else if util.ExecutableExists("pip") {
		pip = "pip"
	} else {
		logrus.Fatal("'pip3' or 'pip' must be installed")
	}

	if err != nil {
		logrus.Fatal(err)
	}

	// install dependencies
	logrus.Info("checking requirements.txt in python project")
	requirementsPath := filepath.Join(tmpDir, "requirements.txt")
	if util.FileExists(requirementsPath) {
		logrus.Info("requirements.txt exists")
	} else {
		logrus.Info("requirements.txt not found. create a new one")

		pipreqs := "pipreqs"
		if !util.ExecutableExists(pipreqs) {
			logrus.Error("'pipreqs' not installed")
			logrus.Fatalf("please run '%s install pipreqs' to install dependency", pip)
		}

		out, err := util.RunCommand(pipreqs, tmpDir, "--savepath", requirementsPath)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info(out)
	}

	logrus.Info("install python libraries")
	out, err := util.RunCommand(pip, "install", "-r", requirementsPath, "-t", tmpDir)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(out)

	// zip
	logrus.Infof("zip python project with dependencies to %s", target)
	err = archive.Zip(tmpDir, target)
	if err != nil {
		logrus.Fatal("Failed to zip python project. ", err)
	}

	logrus.Info(fmt.Sprintf("successfully compressed %s into %s", projDir, target))
}
