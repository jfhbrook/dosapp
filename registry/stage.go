package registry

import (
	"fmt"
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

func (ch *Stage) StagedPackagePath(pkg *Package) string {
	return filepath.Join(
		ch.Config.PackageStageHome,
		fmt.Sprintf("%s-%s-%s",
			pkg.Name,
			pkg.UpstreamVersion,
			pkg.UpstreamReleaseVersion,
		),
	)
}

func (ch *Stage) StagedPackageExists(pkg *Package) bool {
	_, err := os.Stat(ch.StagedPackagePath(pkg))
	return err == nil
}

func (ch *Stage) RemoveStagedPackage(pkg *Package) error {
	return os.RemoveAll(ch.StagedPackagePath(pkg))
}
