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
	"os"
	"path/filepath"

	"github.com/BlueWhaleKo/python-lambda-zipper/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type dockerArgs struct {
	ProjectPath string
	OutputImage string
	BaseImage   string
}

var dockerargs = &dockerArgs{}

var dockerCmd = &cobra.Command{
	Use:   "docker --[flags] [options]",
	Short: "pack python projects into a zip file",
	Long:  "pack python projects into a zip file",
	Run: func(cmd *cobra.Command, args []string) {
		main(cmd, args)
	},
}

func NewDockerCommand() *cobra.Command {
	return dockerCmd
}

func init() {
	// parse args
	dockerCmd.Flags().StringVar(&dockerargs.ProjectPath, "project-path", "", "(required) path to python project")
	dockerCmd.Flags().StringVar(&dockerargs.OutputImage, "output-image", "", "(required) name of output image")
	dockerCmd.Flags().StringVar(&dockerargs.BaseImage, "base-image", "", "(required) name of base image to build from")

	dockerCmd.MarkFlagRequired("project-path")
	dockerCmd.MarkFlagRequired("output-image")
	dockerCmd.MarkFlagRequired("base-image")
}

func validate() error {
	if !util.ExecutableExists("docker") {
		return fmt.Errorf("You need docker installed to run this command")
	}

	if !util.FileExists(filepath.Join(dockerargs.ProjectPath, "__main__.py")) {
		return fmt.Errorf("You need __main__.py at python project root %s as entrypoint", dockerargs.ProjectPath)
	}

	return nil
}

func createDockerfile() *util.Dockerfile {
	// create Dockerfile
	targetDir := "/app"

	d := util.NewDockerfile(dockerargs.BaseImage)
	d.Run("pip install pipreqs")
	d.Add("./", targetDir)
	d.WORKDIR(targetDir)
	d.Run("pipreqs .")
	d.Run("pip install -r ./requirements.txt -t .")

	d.Entrypoint("python __main__.py")

	return d
}

func writeDockerfile(d *util.Dockerfile) (string, error) {
	lines := d.Build()
	logrus.Info("Dockerfile: \n" + lines)

	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}

	err = util.Write(tmpFile.Name(), lines)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil

}

func main(cmd *cobra.Command, args []string) {
	err := validate()
	if err != nil {
		logrus.Fatal(err)
	}

	dfile := createDockerfile()
	dfpath, err := writeDockerfile(dfile)
	if err != nil {
		logrus.Fatal(err)
	}
	defer os.Remove(dfpath)

	out, err := util.RunCommand("docker", "build", "-t", dockerargs.OutputImage, "-f", dfpath, dockerargs.ProjectPath)

	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(out)
}
