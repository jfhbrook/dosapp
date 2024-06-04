package registry

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

func (pkg *Package) Remove() error {
	return os.RemoveAll(pkg.Path())
}

func (pkg *Package) Exists() bool {
	_, err := os.Stat(pkg.Path())
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
