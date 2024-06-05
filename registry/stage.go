package registry

import (
	"os"
	"path/filepath"

	"github.com/jfhbrook/dosapp/config"
)

type Stage struct {
	Config *config.Config
}

func newStage(conf *config.Config) *Stage {
	return &Stage{
		Config: conf,
	}
}

func (ch *Stage) Mkdir() error {
	return os.MkdirAll(ch.Config.PackageStageHome, 0755)
}

func (ch *Stage) Clear() error {
	err := os.RemoveAll(ch.Config.PackageStageHome)
	if err != nil {
		return err
	}
	return ch.Mkdir()
}

func (ch *Stage) StagedPackagePath(name string) string {
	return filepath.Join(ch.Config.PackageStageHome, name)
}

func (ch *Stage) StagedPackageExists(name string) bool {
	_, err := os.Stat(ch.StagedPackagePath(name))
	return err == nil
}

func (ch *Stage) RemoveStagedPackage(name string) error {
	return os.RemoveAll(ch.StagedPackagePath(name))
}
