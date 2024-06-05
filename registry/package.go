package registry

import (
	"archive/tar"
	"compress/gzip"
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
	artifactPath := pkg.LocalArtifactPath()

	if err := pkg.MkArtifactDir(); err != nil {
		return err
	}

	log.Info().Str(
		"name", pkg.Name,
	).Str(
		"url", pkg.URL,
	).Str(
		"path", artifactPath,
	).Msgf("Fetching %s...", pkg.Name)

	resp, err := http.Get(pkg.URL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var f *os.File
	f, err = os.Create(artifactPath)
	defer f.Close()

	if err != nil {
		log.Warn().Err(err).Str("path", artifactPath).Msg("Failed to create artifact")
		return err
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		pkg.Name,
	)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)

	return err
}

func (pkg *Package) Unpack() error {
	artifactPath := pkg.LocalArtifactPath()
	cachedPath := pkg.Cache.CachedPackagePath(pkg.Name)

	// TODO: Ideally I would create an App, but that then creates a Package.
	// I'll need to refactor it to inject the registry. But if we're injecting
	// the registry, maybe I should create it in the Config?
	appPath := filepath.Join(pkg.Config.ConfigHome, "apps", pkg.Name)

	log.Debug().Msgf("Unpacking %s (%s)", pkg.Name, artifactPath)

	artifact, err := os.Open(artifactPath)
	defer artifact.Close()

	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(artifact)

	if err != nil {
		return err
	}

	tar := tar.NewReader(gzipReader)

	log.Debug().Str(
		"staging", pkg.Cache.CachedPackagePath(pkg.Name),
	).Msg("Cleaning up staged package")

	if err := pkg.Cache.RemoveCachedPackage(pkg.Name); err != nil {
		return err
	}

	log.Info().Str(
		"package", pkg.Name,
	).Str(
		"source", artifactPath,
	).Str(
		"destination", cachedPath,
	).Msg("Unpacking artifact...")

	for {
		hdr, err := tar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filename := filepath.Join(cachedPath, hdr.Name)

		log.Debug().Str("file", filename).Msgf("Unpacking %s...", filename)
		f, err := os.Create(filename)
		defer f.Close()

		if err != nil {
			return err
		}

		bar := progressbar.DefaultBytes(-1, filename)

		if _, err := io.Copy(io.MultiWriter(f, bar), tar); err != nil {
			return err
		}
	}

	log.Debug().Str(
		"package", pkg.Name,
	).Str(
		"path", artifactPath,
	).Msg("Unpacking complete")

	log.Info().Str(
		"package", pkg.Name,
	).Str(
		"source", artifactPath,
	).Str(
		"destination", cachedPath,
	).Msg("Installing package...")

	err = os.RemoveAll(appPath)

	if err != nil {
		return err
	}

	// NOTE: There are some limitations to this approach - for instance, rename
	// can't copy between drives.
	err = os.Rename(cachedPath, appPath)

	if err != nil {
		return err
	}

	log.Info().Str(
		"package", pkg.Name,
	).Str(
		"source", artifactPath,
	).Str(
		"destination", cachedPath,
	).Msg("Installation complete")

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
