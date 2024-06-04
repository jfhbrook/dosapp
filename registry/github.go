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

func NewGitHubRelease(name string, version string, releaseVersion string, url string) (*GitHubRelease, error) {
	ver, err := semver.NewVersion(version)

	if err != nil {
		return nil, err
	}

	var relVer *semver.Version
	relVer, err = semver.NewVersion(releaseVersion)

	if err != nil {
		return nil, err
	}

	return &GitHubRelease{
		Name:           name,
		Version:        ver,
		ReleaseVersion: relVer,
		Url:            url,
	}, nil
}

// TODO: This is the url for the HTML list of assets, not the actual asset
// we want. That's going to be on Assets[0].Url.
// See: https://github.com/google/go-github/blob/446a2968d80b0adc88c54a1c1a17a05b25e83ea7/github/repos_releases.go#L68
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

func (reg *GitHubRegistry) GitHubRelease(name string) (*GitHubRelease, error) {
	releases, _, err := reg.client.Repositories.ListReleases(
		context.Background(),
		reg.user,
		reg.repository,
		nil,
	)

	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		// TODO: A lot of this is factory stuff. Can/should I push it into
		// GitHubRelease?
		tag := *release.TagName

		log.Info().Msgf("Tag: %s", tag)

		sp := strings.Split(tag, "-")
		nameSp, version, releaseVersion := sp[:len(sp)-2], sp[len(sp)-2], sp[len(sp)-1]
		releaseName := strings.Join(nameSp, "-")

		url := *release.AssetsURL

		if releaseName == name {
			log.Debug().Msgf("Package %s found", name)
			return NewGitHubRelease(
				name,
				version,
				releaseVersion,
				url,
			)
		}
	}
	log.Debug().Msgf("Package %s not found", name)
	return nil, nil
}

func (reg *GitHubRegistry) PackageURL(name string) (string, error) {
	release, err := reg.GitHubRelease(name)

	if err != nil {
		return "", err
	}

	log.Debug().Msgf("Name: %s", release.Name)
	log.Debug().Msgf("Version: %s", release.Version)
	log.Debug().Msgf("GitHubRelease Version: %s", release.ReleaseVersion)
	log.Debug().Msgf("GitHubRelease URL: %s", release.Url)

	return release.Url, nil
}
