package packages

import (
	"os"
	"path/filepath"

	"github.com/jfhbrook/dosapp/config"
)

type Artifact struct {
	Name   string
	Config *config.Config
}

func NewArtifact(conf *config.Config, name string) *Artifact {
	return &Artifact{
		Name:   name,
		Config: conf,
	}
}

func (artifact *Artifact) Path() string {
	return filepath.Join(artifact.Config.ArtifactHome, artifact.Name+".tar.gz")
}

func (artifact *Artifact) Remove() error {
	if err := os.RemoveAll(artifact.Path()); err != nil {
		return err
	}
	return os.RemoveAll(artifact.Path())
}

func (artifact *Artifact) Exists() bool {
	_, err := os.Stat(artifact.Path())
	return err == nil
}

func (artifact *Artifact) Url() error {
	return nil
}
