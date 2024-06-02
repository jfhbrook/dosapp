package packages

import (
	"os"
	"path/filepath"

	"github.com/jfhbrook/dosapp/config"
)

type Package struct {
	Name   string
	Config *config.Config
}

func NewPackage(conf *config.Config, name string) *Package {
	return &Package{
		Name:   name,
		Config: conf,
	}
}

func (pkg *Package) Path() string {
	return filepath.Join(pkg.Config.PackageHome, pkg.Name)
}

// TODO: Not sure this is where I want to put tarballs
func (pkg *Package) ArtifactPath() string {
	return filepath.Join(pkg.Config.DownloadHome, "packages", pkg.Name+".tar.gz")
}

func (pkg *Package) Remove() error {
	if err := os.RemoveAll(pkg.ArtifactPath()); err != nil {
		return err
	}
	return os.RemoveAll(pkg.Path())
}

func (pkg *Package) Exists() bool {
	_, err := os.Stat(pkg.Path())
	return err == nil
}

func (pkg *Package) ArtifactExists() bool {
	_, err := os.Stat(pkg.ArtifactPath())
	return err == nil
}

func (pkg *Package) EnvFileTemplatePath() string {
	return filepath.Join(pkg.Path(), "dosapp.env.tmpl")
}

func (pkg *Package) TaskFilePath() string {
	return filepath.Join(pkg.Path(), "Taskfile.yml")
}

func (pkg *Package) DocsPath() string {
	return filepath.Join(pkg.Path(), "README.md")
}
