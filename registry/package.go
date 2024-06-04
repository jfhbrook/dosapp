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
	Cache                  *Cache
	Config                 *config.Config
}

func localPackagePath(name string, conf *config.Config) string {
	return filepath.Join(conf.PackageHome, name)
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
	ch *Cache,
	conf *config.Config,
) *Package {
	var pkg *Package
	var err error

	pkg, err = fromPackageFile(name, conf)
	pkg.Cache = ch

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

func (pkg *Package) LocalArtifactPath() string {
	return filepath.Join(pkg.Config.ArtifactHome, pkg.Name+".tar.gz")
}

func (pkg *Package) LocalArtifactExists() bool {
	_, err := os.Stat(pkg.LocalArtifactPath())
	return err == nil
}

func (pkg *Package) RemoveLocalArtifact() error {
	return os.RemoveAll(pkg.LocalArtifactPath())
}

func (pkg *Package) LocalPackagePath() string {
	return filepath.Join(pkg.Config.PackageHome, pkg.Name)
}

func (pkg *Package) LocalPackageExists() bool {
	_, err := os.Stat(pkg.LocalPackagePath())
	return err == nil
}

func (pkg *Package) RemoveLocalPackage() error {
	return os.RemoveAll(pkg.LocalPackagePath())
}

func (pkg *Package) CachedPackagePath() string {
	return pkg.Cache.CachedPackagePath(pkg.Name)
}

func (pkg *Package) CachedPackageExists() bool {
	return pkg.Cache.CachedPackageExists(pkg.Name)
}

func (pkg *Package) RemoveCachedPackage() error {
	return pkg.Cache.RemoveCachedPackage(pkg.Name)
}

// TODO: Download package from URL
func (pkg *Package) Download() error {
	return nil
}

// TODO: Unpack package into cache, then move into package home
func (pkg *Package) Unpack() error {
	return nil
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
