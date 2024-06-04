package registry

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/rs/zerolog/log"

	"github.com/jfhbrook/dosapp/config"
)

type RegistryError struct {
	message string
}

func (e *RegistryError) Error() string {
	return e.message
}

type Registry interface {
	PackageURL(name string) (string, error)
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

func (reg *GitHubRegistry) Release(name string) (*Release, error) {
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
		// Release?
		tag := *release.TagName

		log.Info().Msgf("Tag: %s", tag)

		sp := strings.Split(tag, "-")
		nameSp, version, releaseVersion := sp[:len(sp)-2], sp[len(sp)-2], sp[len(sp)-1]
		releaseName := strings.Join(nameSp, "-")

		url := *release.AssetsURL

		if releaseName == name {
			log.Debug().Msgf("Package %s found", name)
			return NewRelease(
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
	release, err := reg.Release(name)

	if err != nil {
		return "", err
	}

	log.Debug().Msgf("Name: %s", release.Name)
	log.Debug().Msgf("Version: %s", release.Version)
	log.Debug().Msgf("Release Version: %s", release.ReleaseVersion)
	log.Debug().Msgf("Release URL: %s", release.Url)

	return release.Url, nil
}

func NewRegistry(conf *config.Config) (Registry, error) {
	u, err := url.Parse(conf.Registry)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "github" {
		return nil, &RegistryError{fmt.Sprintf("Forge %s is unsupported.", u.Scheme)}
	}

	return newGitHubRegistry(conf, u)
}
