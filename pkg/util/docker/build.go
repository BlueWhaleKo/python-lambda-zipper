package util

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func BuildImageFromPath(c *client.Client, path string, opt types.ImageBuildOptions) (types.ImageBuildResponse, error) {

	buildCtx, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		return types.ImageBuildResponse{}, err
	}

	buildOpts := types.ImageBuildOptions{
		//Dockerfile: dockerfileName,
	}

	return c.ImageBuild(context.TODO(), buildCtx, buildOpts)
}

func BuildImage(c *client.Client, dockerfile io.Reader, opt types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	readall, err := ioutil.ReadAll(dockerfile)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}

	err = tw.WriteHeader(&tar.Header{
		Name: "Dockerfile",
		Size: int64(len(readall)),
	})
	if err != nil {
		return types.ImageBuildResponse{}, err
	}

	_, err = tw.Write(readall)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}

	tarReader := bytes.NewReader(buf.Bytes())
	res, err := c.ImageBuild(context.Background(), tarReader, opt)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}
	return res, nil
}
