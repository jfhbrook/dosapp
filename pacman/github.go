package pacman

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v62/github"
	"github.com/rs/zerolog/log"

	"github.com/jfhbrook/dosapp/config"
)

type GitHubRelease struct {
	Name           string
	Version        *semver.Version
	ReleaseVersion *semver.Version
	URL            string
}

func newGitHubRelease(ghRelease *github.RepositoryRelease) (*GitHubRelease, error) {
	tag := *ghRelease.TagName

	log.Debug().Str("tag", tag).Msgf("Loading release: %s", tag)

	sp := strings.Split(tag, "-")
	nameSp, ver, relVer := sp[:len(sp)-2], sp[len(sp)-2], sp[len(sp)-1]
	releaseName := strings.Join(nameSp, "-")

	url := *ghRelease.Assets[0].BrowserDownloadURL

	version, err := semver.NewVersion(ver)

	if err != nil {
		return nil, err
	}

	var releaseVersion *semver.Version
	releaseVersion, err = semver.NewVersion(relVer)

	if err != nil {
		return nil, err
	}

	release := &GitHubRelease{
		Name:           releaseName,
		Version:        version,
		ReleaseVersion: releaseVersion,
		URL:            url,
	}

	return release, nil
}

type GitHubRegistry struct {
	stage      *Stage
	Config     *config.Config
	client     *github.Client
	user       string
	repository string
}

func newGitHubRegistry(conf *config.Config, u *url.URL) Registry {
	user := u.Host
	re := regexp.MustCompile(`^/`)
	repo := re.ReplaceAllString(u.Path, "")

	ch := newStage(conf)

	return &GitHubRegistry{
		stage:      ch,
		Config:     conf,
		client:     github.NewClient(nil),
		user:       user,
		repository: repo,
	}
}

func (reg *GitHubRegistry) Stage() *Stage {
	return reg.stage
}

func (reg *GitHubRegistry) Release(name string) (*GitHubRelease, error) {
	releases, _, err := reg.client.Repositories.ListReleases(
		context.Background(),
		reg.user,
		reg.repository,
		nil,
	)

	if err != nil {
		return nil, err
	}

	for _, ghRelease := range releases {
		release, err := newGitHubRelease(ghRelease)

		if err != nil {
			return nil, err
		}

		if release.Name == name {
			return release, nil
		}
	}
	return nil, nil
}

func (reg *GitHubRegistry) FindPackage(name string) *Package {
	release, err := reg.Release(name)

	if err != nil {
		log.Debug().Err(err).Msgf("Upstream package %s not found", name)

		return NewPackage(
			name,
			nil,
			nil,
			"",
			reg.Stage(),
			reg.Config,
		)
	} else {
		log.Debug().Str(
			"name", name,
		).Str(
			"version", release.Version.String(),
		).Str(
			"releaseVersion", release.ReleaseVersion.String(),
		).Str(
			"url", release.URL,
		).Msgf("Upstream package %s found", name)

		return NewPackage(
			name,
			release.Version,
			release.ReleaseVersion,
			release.URL,
			reg.Stage(),
			reg.Config,
		)
	}
}
