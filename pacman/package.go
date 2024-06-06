package pacman

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"

	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/manifest"
)

type Package struct {
	Name                   string
	remotePackageExists    bool
	LocalVersion           *semver.Version
	LocalReleaseVersion    *semver.Version
	UpstreamVersion        *semver.Version
	UpstreamReleaseVersion *semver.Version
	URL                    string
	Stage                  *Stage
	Config                 *config.Config
}

func localPackagePath(name string, conf *config.Config) string {
	return filepath.Join(conf.PackageHome, name)
}

// TODO: This call signature is great for a remote package, but many of my
// use cases are for packages only installed locally - and those are having
// to inject a bunch of null values.
func NewPackage(
	name string,
	version *semver.Version,
	releaseVersion *semver.Version,
	url string,
	ch *Stage,
	conf *config.Config,
) *Package {
	packagePath := filepath.Join(localPackagePath(name, conf), "Package.yml")

	pkg := &Package{}
	pkg.Name = name
	pkg.Config = conf

	pkg.Stage = ch
	pkg.remotePackageExists = version != nil && releaseVersion != nil && url != ""
	pkg.UpstreamVersion = version
	pkg.UpstreamReleaseVersion = releaseVersion
	pkg.URL = url

	m, err := manifest.FromFile(packagePath)

	if err != nil {
		log.Debug().Str("path", packagePath).Err(err).Msg("Installed package not found")
	} else {
		pkg.LocalVersion = m.Version
		pkg.LocalReleaseVersion = m.ReleaseVersion
	}
	return pkg
}

func (pkg *Package) RemotePackageExists() bool {
	return pkg.remotePackageExists
}

func (pkg *Package) LocalArtifactPath() string {
	return filepath.Join(
		pkg.Config.ArtifactHome,
		fmt.Sprintf(
			"%s-%s-%s.tar.gz",
			pkg.Name,
			pkg.UpstreamVersion.String(),
			pkg.UpstreamReleaseVersion.String(),
		),
	)
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

func (pkg *Package) MkPackageHome() error {
	return os.MkdirAll(pkg.Config.PackageHome, 0755)
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

func (pkg *Package) StagedPackagePath() string {
	return pkg.Stage.StagedPackagePath(pkg)
}

func (pkg *Package) StagedPackageExists() bool {
	return pkg.Stage.StagedPackageExists(pkg)
}

func (pkg *Package) RemoveStagedPackage() error {
	return pkg.Stage.RemoveStagedPackage(pkg)
}

func (pkg *Package) HasUpdate() bool {
	if pkg.LocalVersion == nil {
		return true
	}

	if pkg.LocalVersion.LessThan(pkg.UpstreamVersion) {
		return true
	}

	if pkg.LocalVersion.Equal(pkg.UpstreamVersion) && pkg.LocalReleaseVersion.LessThan(pkg.UpstreamReleaseVersion) {
		return true
	}

	return false
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
	packagePath := pkg.LocalPackagePath()

	log.Debug().Msgf("Unpacking %s (%s)", pkg.Name, artifactPath)

	artifact, err := os.Open(artifactPath)
	defer artifact.Close()

	if err != nil {
		log.Info().Msg("lol")
		return err
	}

	gzipReader, err := gzip.NewReader(artifact)

	if err != nil {
		return err
	}

	tar := tar.NewReader(gzipReader)

	log.Debug().Str(
		"staging", pkg.Config.PackageStageHome,
	).Msg("Setting up staging area...")

	if err := pkg.Stage.Mkdir(); err != nil {
		return err
	}

	if err := pkg.Stage.RemoveStagedPackage(pkg); err != nil {
		return err
	}

	log.Info().Str(
		"package", pkg.Name,
	).Str(
		"source", artifactPath,
	).Msg("Unpacking artifact...")

	for {
		hdr, err := tar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filename := filepath.Join(pkg.Config.PackageStageHome, hdr.Name)

		if hdr.FileInfo().IsDir() {
			log.Debug().Str(
				"directory", filename,
			).Msgf("Creating %s...", filename)

			os.Mkdir(filename, 0755)
			continue
		}

		log.Debug().Str("file", filename).Msg("Unpacking file...")
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

	stagedPath := pkg.StagedPackagePath()

	log.Debug().Str(
		"package", pkg.Name,
	).Str(
		"path", artifactPath,
	).Msg("Unpacking complete")

	log.Debug().Str(
		"DOSAPP_PACKAGE_HOME", pkg.Config.PackageHome,
	).Msg("Setting up package home...")

	err = pkg.MkPackageHome()

	if err != nil {
		return err
	}

	log.Info().Str(
		"package", pkg.Name,
	).Str(
		"source", artifactPath,
	).Str(
		"destination", packagePath,
	).Msg("Installing package...")

	err = os.RemoveAll(packagePath)

	if err != nil {
		return err
	}

	// NOTE: There are some limitations to this approach - for instance, rename
	// can't copy between drives.
	err = os.Rename(stagedPath, packagePath)

	if err != nil {
		return err
	}

	log.Info().Str(
		"package", pkg.Name,
	).Str(
		"source", artifactPath,
	).Str(
		"destination", packagePath,
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
