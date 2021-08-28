package util

import (
	"fmt"
	"os"
	"strings"
)

type Dockerfile struct {
	inner []string
}

func NewDockerfile() *Dockerfile {
	return &Dockerfile{}
}

func (d *Dockerfile) From(from string) *Dockerfile {
	d.inner = append(d.inner, fmt.Sprintf("FROM %s", from))
	return d
}

func (d *Dockerfile) Entrypoint(entrypoint ...string) *Dockerfile {
	s := strings.Join(entrypoint, " ")
	d.inner = append(d.inner, fmt.Sprintf("ENTRYPOINT %s", s))
	return d
}

func (d *Dockerfile) Cmd(cmd ...string) *Dockerfile {
	s := strings.Join(cmd, " ")
	d.inner = append(d.inner, fmt.Sprintf("CMD %s", s))
	return d
}

func (d *Dockerfile) WORKDIR(dir string) *Dockerfile {
	d.inner = append(d.inner, fmt.Sprintf("WORKDIR %s", dir))
	return d
}

func (d *Dockerfile) Run(cmd ...string) *Dockerfile {
	s := strings.Join(cmd, " ")
	d.inner = append(d.inner, fmt.Sprintf("RUN %s", s))
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
	return strings.Join(d.inner, "\n")
}

func (d *Dockerfile) WriteTo(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, d.Build())
	if err != nil {
		return err
	}

	return nil
}
