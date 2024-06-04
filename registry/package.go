package registry

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
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
	pkg.Name = name
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

func (pkg *Package) MkArtifactDir() error {
	return os.MkdirAll(pkg.Config.ArtifactHome, 0755)
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

func (pkg *Package) Fetch() error {
	if err := pkg.MkArtifactDir(); err != nil {
		return err
	}

	resp, err := http.Get(pkg.URL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var f *os.File
	f, err = os.Create(pkg.LocalArtifactPath())
	defer f.Close()

	if err != nil {
		log.Warn().Err(err).Msg("Failed to create file")
		return err
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading "+pkg.Name+"...",
	)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)

	return err
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
