package registry

import (
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/jfhbrook/dosapp/config"
)

type Package struct {
	Name                   string `yaml:"name"`
	remotePackageExists    bool
	LocalVersion           *semver.Version `yaml:"version"`
	LocalReleaseVersion    *semver.Version `yaml:"release_version"`
	UpstreamVersion        *semver.Version
	UpstreamReleaseVersion *semver.Version
	URL                    string
	Config                 *config.Config
}

func localPackagePath(name string, conf *config.Config) string {
	return filepath.Join(conf.PackageHome, name)
}

func localArtifactPath(name string, conf *config.Config) string {
	return filepath.Join(conf.ArtifactHome, name+".tar.gz")
}

func fromPackageFile(name string, conf *config.Config) (*Package, error) {
	var file []byte
	var err error
	file, err = os.ReadFile(filepath.Join(localPackagePath(name, conf), "Package.yml"))

	pkg := &Package{}
	pkg.Config = conf

	if err != nil {
		return pkg, err
	} else {
		err = yaml.Unmarshal(file, pkg)
		if err != nil {
			return pkg, err
		}
	}

	return pkg, nil
}

func newPackage(
	name string,
	version *semver.Version,
	releaseVersion *semver.Version,
	url string,
	conf *config.Config,
) *Package {
	var pkg *Package
	var err error

	pkg, err = fromPackageFile(name, conf)

	if err != nil {
		log.Debug().Err(err).Msg("Package not found locally")
	}

	pkg.remotePackageExists = version != nil && releaseVersion != nil && url != ""
	pkg.UpstreamVersion = version
	pkg.UpstreamReleaseVersion = releaseVersion
	pkg.URL = url

	return pkg
}

func (pkg *Package) RemotePackageExists() bool {
	return pkg.remotePackageExists
}

func (pkg *Package) LocalArtifactExists() bool {
	_, err := os.Stat(localArtifactPath(pkg.Name, pkg.Config))
	return err == nil
}

func (pkg *Package) LocalPackagePath() string {
	return filepath.Join(pkg.Config.PackageHome, pkg.Name)
}

func (pkg *Package) Remove() error {
	return os.RemoveAll(pkg.LocalPackagePath())
}

func (pkg *Package) Exists() bool {
	_, err := os.Stat(pkg.LocalPackagePath())
	return err == nil
}

func (pkg *Package) EnvFileTemplatePath() string {
	return filepath.Join(pkg.LocalPackagePath(), "dosapp.env.tmpl")
}

func (pkg *Package) TaskFilePath() string {
	return filepath.Join(pkg.LocalPackagePath(), "Taskfile.yml")
}

func (pkg *Package) DocsPath() string {
	return filepath.Join(pkg.LocalPackagePath(), "README.md")
}
