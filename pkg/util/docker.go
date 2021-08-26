package util

import (
	"fmt"
	"strings"
)

type Dockerfile struct {
	image      string
	entrypoint string
	cmd        string
	inner      []string
}

func NewDockerfile(image string) *Dockerfile {
	inner := []string{"FROM " + image}
	return &Dockerfile{
		image: image,
		inner: inner,
	}
}

func (d *Dockerfile) Entrypoint(cmd string) *Dockerfile {
	d.entrypoint = fmt.Sprintf("ENTRYPOINT %s", cmd)
	return d
}

func (d *Dockerfile) Cmd(cmd string) *Dockerfile {
	d.cmd = fmt.Sprintf("CMD %s", cmd)
	return d
}

func (d *Dockerfile) WORKDIR(dir string) *Dockerfile {
	d.inner = append(d.inner, fmt.Sprintf("WORKDIR %s", dir))
	return d
}

func (d *Dockerfile) Run(cmd string) *Dockerfile {
	d.inner = append(d.inner, fmt.Sprintf("RUN %s", cmd))
	return d
}
func (d *Dockerfile) Add(source, target string) *Dockerfile {
	d.inner = append(d.inner, fmt.Sprintf("ADD %s %s", source, target))
	return d
}

func (d *Dockerfile) Copy(source, target string) *Dockerfile {
	d.inner = append(d.inner, fmt.Sprintf("COPY %s %s", source, target))
	return d
}

func (d *Dockerfile) Build() string {
	lines := d.inner
	lines = append(lines, d.entrypoint)
	lines = append(lines, d.cmd)

	return strings.Join(lines, "\n")
}
