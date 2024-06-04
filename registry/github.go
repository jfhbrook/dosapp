package registry

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
	Url            string
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
		Url:            url,
	}

	return release, nil
}

func (release *GitHubRelease) URL() string {
	return release.Url
}

type GitHubRegistry struct {
	Config     *config.Config
	client     *github.Client
	user       string
	repository string
}

func newGitHubRegistry(conf *config.Config, u *url.URL) (Registry, error) {
	user := u.Host
	re := regexp.MustCompile(`^/`)
	repo := re.ReplaceAllString(u.Path, "")

	return &GitHubRegistry{
		Config:     conf,
		client:     github.NewClient(nil),
		user:       user,
		repository: repo,
	}, nil
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
			log.Debug().Str(
				"name", name,
			).Str(
				"version", release.Version.String(),
			).Str(
				"releaseVersion", release.ReleaseVersion.String(),
			).Str(
				"url", release.Url,
			).Msgf("Package %s found", name)

			return release, nil
		}
	}
	log.Debug().Msgf("Package %s not found", name)
	return nil, nil
}

func (reg *GitHubRegistry) PackageURL(name string) (string, error) {
	release, err := reg.Release(name)

	if err != nil {
		return "", err
	}

	return release.Url, nil
}
